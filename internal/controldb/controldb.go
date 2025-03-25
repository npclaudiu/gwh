package controldb

import (
	"database/sql"
	"fmt"
	"regexp"
	"slices"
	"strconv"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

const (
	RegKeySchemaVersion string = "schema_version"
)

type ControlDatabase struct {
	db *sql.DB
}

func Open(path string) (*ControlDatabase, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("gwh: failed to open control database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gwh: failed to ping control database: %w", err)
	}

	cdb := &ControlDatabase{db: db}

	current, err := cdb.loadSchemaVersion()

	if err != nil {
		return nil, fmt.Errorf("gwh: failed to load current control database schema version: %w", err)
	}

	embedded, err := cdb.listEmbeddedSchemaVersions()

	if err != nil {
		return nil, fmt.Errorf("gwh: failed to list schema versions: %w", err)
	}

	highest := slices.Max(embedded)

	if current < highest {
		cdb.migrateSchema(current, highest)
	}

	return cdb, nil
}

func (c *ControlDatabase) Close() error {
	err := c.db.Close()

	if err != nil {
		return fmt.Errorf("gwh: failed to close control database: %w", err)
	}

	return nil
}

//#region Primitives

func (c *ControlDatabase) queryScalar(q string, result any) error {
	rows, err := c.db.Query(q)

	if err != nil {
		return fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(result); err != nil {
			return fmt.Errorf("failed to scan count: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to iterate through rows")
}

//#endregion

//#region Schema

func (c *ControlDatabase) loadSchemaVersion() (int, error) {
	hasRegistryTable, err := c.hasRegistryTable()

	if err != nil {
		return 0, fmt.Errorf("failed to check existance of the registry table: %w", err)
	}

	if !hasRegistryTable {
		return 0, nil
	}

	schemaVersionRaw, err := c.loadRegistryValue(RegKeySchemaVersion)

	if err != nil {
		return 0, fmt.Errorf("failed to load schema version: %w", err)
	}

	schemaVersion := 0

	if schemaVersionRaw != "" {
		schemaVersion, err = strconv.Atoi(schemaVersionRaw)

		if err != nil {
			return 0, fmt.Errorf("failed to parse schema version as int: %w", err)
		}
	}

	return schemaVersion, nil
}

func (c *ControlDatabase) listEmbeddedSchemaVersions() ([]int, error) {
	schemas, err := embeddedSchemas.ReadDir("schema")

	if err != nil {
		return nil, fmt.Errorf("gwh: failed to read schema directory: %w", err)
	}

	versions := make([]int, len(schemas))
	re, err := regexp.Compile(`v(\d+)\.sql`)

	if err != nil {
		return nil, fmt.Errorf("failed to compile regex: %w", err)
	}

	for i, schema := range schemas {
		name := schema.Name()
		match := re.FindStringSubmatch(name)
		version, err := strconv.Atoi(match[1])

		if err != nil {
			return nil, fmt.Errorf("failed to parse schema version: %w", err)
		}

		versions[i] = version
	}

	slices.Sort(versions)

	return versions, nil
}

func (c *ControlDatabase) migrateSchema(from, to int) error {
	if from >= to {
		return fmt.Errorf("invalid migration range: %d -> %d", from, to)
	}

	tx, err := c.db.Begin()

	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	for i := from + 1; i <= to; i++ {
		q := fmt.Sprintf("schema/v%d.sql", i)
		schema, err := embeddedSchemas.ReadFile(q)

		if err != nil {
			return fmt.Errorf("failed to read schema file: %w", err)
		}

		if _, err := tx.Exec(string(schema)); err != nil {
			return fmt.Errorf("failed to execute schema migration: %w", err)
		}
	}

	q := `
		insert into gwh_registry (key, value)
			values (?, ?)
			on conflict (key) do update set value = excluded.value;
	`

	if _, err := tx.Exec(q, RegKeySchemaVersion, to); err != nil {
		return fmt.Errorf("failed to update schema version: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

//#endregion

//#region Registry

func (c *ControlDatabase) hasRegistryTable() (bool, error) {
	q := `
		select count(*) as count
			from sqlite_master
			where type='table' and name='gwh_registry';
	`

	var count int
	err := c.queryScalar(q, &count)

	if err != nil {
		return false, fmt.Errorf("query failed: %w", err)
	}

	return count > 0, nil
}

func (c *ControlDatabase) loadRegistryValue(key string) (string, error) {
	q := `
		select value
			from gwh_registry
			where key = ?;
	`

	rows, err := c.db.Query(q, key)

	if err != nil {
		return "", fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	if rows.Next() {
		var value string

		if err := rows.Scan(&value); err != nil {
			return "", fmt.Errorf("failed to scan value: %w", err)
		}

		return value, nil
	}

	return "", nil
}

//#endregion

//#region Linking

func (c *ControlDatabase) LinkRepository(name, path string) error {
	q := `
		insert into gwh_git_repositories (name, path)
			values (?, ?);
	`

	if _, err := c.db.Exec(q, name, path); err != nil {
		sqliteErr, ok := err.(sqlite3.Error)

		// Ignore unique constraint errors. Handling this case in code rather
		// than using `on conflict (...) do nothing` in SQL for flexibility.
		if ok && sqliteErr.Code != sqlite3.ErrConstraint {
			return fmt.Errorf("failed to link repository: %w", err)
		}
	}

	// TODO(npclaudiu): Add support for recursively linking submodules of
	// the repository being linked.

	return nil
}

//#endregion
