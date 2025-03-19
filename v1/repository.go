package gwh

import (
	"context"
	"database/sql"
	"os"
	"path"
)

type AddRepositoryOptions struct {
	Prefix string
	Path   string
	Sync   bool
}

func AddRepository(_ context.Context, options *AddRepositoryOptions) error {
	// Check if the working directory is valid.
	//
	prefix, err := validatePrefix(options.Prefix)

	if err != nil {
		return err
	}

	// Check the warehouse directory.
	//
	whDir := path.Join(prefix, WarehouseDirName)
	whDirStat, err := os.Stat(whDir)

	if err != nil {
		return err
	}

	if !whDirStat.IsDir() {
		return ErrWarehouseDirIsInvalid
	}

	// Read manifest.
	//
	manifestPath := path.Join(whDir, ManifestFileName)
	manifest, err := ReadManifest(manifestPath)

	if err != nil {
		return err
	}

	// Add repository to manifest.
	//
	localRepositoryID, err := manifest.AddRepository(ManifestRepositoryKindGit, options.Path)

	if err != nil {
		return err
	}

	// Create repository directory and database.
	//
	repositoryDir := path.Join(whDir, "repositories", localRepositoryID)

	if err := os.MkdirAll(repositoryDir, 0755); err != nil {
		return err
	}

	db, err := sql.Open("duckdb", path.Join(repositoryDir, AnalyticsDatabaseFileName))

	if err != nil {
		return err
	}

	defer db.Close()

	// Sync repository, if requested.
	//
	if options.Sync {
	}

	// Write manifest.
	//
	if err := WriteManifest(manifest, manifestPath); err != nil {
		return err
	}

	return nil
}
