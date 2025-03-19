package gwh

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"

	"database/sql"

	"github.com/charmbracelet/log"
	_ "github.com/marcboeker/go-duckdb"
)

type InitOptions struct {
	Prefix string
}

func Init(_ context.Context, options *InitOptions) error {
	// Check if the working directory is valid.
	//
	prefix, err := validatePrefix(options.Prefix)

	if err != nil {
		return err
	}

	// Create the warehouse directory.
	//
	whDir := path.Join(prefix, WarehouseDirName)
	whDirStat, err := os.Stat(whDir)
	checkWhStat := true

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			checkWhStat = false
		} else {
			log.Error("stat failed", "error", err)
			return ErrWarehouseDirInitFailed
		}
	}

	if checkWhStat && whDirStat.IsDir() {
		return ErrWarehouseDirExists
	}

	if err := os.Mkdir(whDir, 0755); err != nil {
		return ErrWarehouseDirInitFailed
	}

	// Create the manifest file.
	//
	manifest := NewManifest(&ManifestOptions{
		Version: ManifestVersionCurrent,
	})

	if err = WriteManifest(manifest, path.Join(whDir, ManifestFileName)); err != nil {
		return err
	}

	// Create the database.
	//
	// This will need to be revisited, as the current implementation focuses
	// on supporting only one Git repository with no submodules.
	//
	db, err := sql.Open("duckdb", path.Join(whDir, AnalyticsDatabaseFileName))
	if err != nil {
		return err
	}
	defer db.Close()

	// TODO: Bundle duckpgq extension.
	_, err = db.Exec(`install duckpgq from community;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`load duckpgq;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		create table if not exists git_commits (
			id varchar,
			message varchar
		)`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		create table if not exists git_commits_parents (
			commit_id varchar,
			parent_commit_id varchar
		)`,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		create property graph git_commits_pg vertex tables (git_commits)
		edge tables (
			git_commits_parents
				source key (commit_id) references git_commits (id)
				destination key (parent_commit_id) references git_commits (id)
				label parent
  		);
	`)

	if err != nil {
		return err
	}

	return nil
}
