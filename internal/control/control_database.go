package control

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type ControlDatabase struct {
	db *sql.DB
}

func OpenControlDatabase(path string) (*ControlDatabase, error) {
	db, err := sql.Open("sqlite3", path)

	if err != nil {
		return nil, fmt.Errorf("gwh: failed to open control database: %w", err)
	}

	return &ControlDatabase{
		db: db,
	}, nil
}

func (c *ControlDatabase) Close() error {
	return c.db.Close()
}
