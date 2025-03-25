package gwh

import (
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/npclaudiu/gwh/internal/controldb"
	"github.com/npclaudiu/gwh/internal/gitops"
	"github.com/npclaudiu/gwh/internal/gwhlayout"
)

// Warehouse represents a handle that manages a Git warehouse.
type Warehouse struct {
	layout          *gwhlayout.WarehouseLayout
	controlDatabase *controldb.ControlDatabase
}

// Open opens a warehouse at the specified location. If the warehouse does not
// exist, it will be created.
func Open(location string) (*Warehouse, error) {
	// Make sure we have a valid file layout.
	//
	layout, err := gwhlayout.New(location)

	if err != nil {
		return nil, fmt.Errorf("gwh: %w", err)
	}

	// Open control database.
	//
	controlDatabaseFile := layout.GetKnownPath(gwhlayout.ControlDatabaseFile)
	controlDatabase, err := controldb.Open(controlDatabaseFile)

	if err != nil {
		return nil, fmt.Errorf("gwh: failed to open control database: %w", err)
	}

	return &Warehouse{
		layout:          layout,
		controlDatabase: controlDatabase,
	}, nil
}

// Close closes the warehouse handle.
func (w *Warehouse) Close() {
	w.controlDatabase.Close()
}

// LinkRepository links a Git repository to the warehouse.
func (w *Warehouse) LinkRepository(name, repositoryPath string) error {
	// Validate repository name.
	//
	if !w.layout.IsNameValid(name) {
		return fmt.Errorf("gwh: invalid repository name")
	}

	// Open repository.
	//
	path := w.layout.ResolvePath(repositoryPath)
	r, err := gitops.OpenRepository(path)

	if err != nil {
		return fmt.Errorf("gwh: failed to open repository for linking: %w", err)
	}

	storage, ok := r.Storer.(*filesystem.Storage)

	if !ok {
		return fmt.Errorf("gwh: repository storage is not a supported file system")
	}

	// Link repository using its relative path to the warehouse directory.
	//
	dotgit := storage.Filesystem().Root()
	dotgwh := w.layout.GetKnownPath(gwhlayout.WarehouseDirectory)

	path, err = filepath.Rel(dotgwh, dotgit)

	if err != nil {
		return fmt.Errorf("gwh: failed to calculate relative path to repository: %w", err)
	}

	if err := w.controlDatabase.LinkRepository(name, path); err != nil {
		return fmt.Errorf("gwh: failed to link repository: %w", err)
	}

	return nil
}

// SyncRepository pulls data from a linked repository into the warehouse.
func (w *Warehouse) SyncRepository(name string) error {
	rl, err := w.controlDatabase.GetRepositoryLink(name)

	if err != nil {
		return fmt.Errorf("gwh: failed to get repository link: %w", err)
	}

	path := w.layout.ResolvePath(rl.Path)

	repository, err := gitops.OpenRepository(path)

	if err != nil {
		return fmt.Errorf("gwh: failed to open repository for syncing: %w", err)
	}

	branches, err := repository.Branches()

	if err != nil {
		return fmt.Errorf("gwh: failed to get repository branches: %w", err)
	}

	for {
		branch, err := branches.Next()

		if err != nil {
			break
		}

		commits, err := repository.Log(&git.LogOptions{
			From: branch.Hash(),
		})

		if err != nil {
			return fmt.Errorf("gwh: failed to get branch commits: %w", err)
		}

		fmt.Println("Branch:", branch.Name())

		for {
			commit, err := commits.Next()

			if err != nil {
				break
			}

			fmt.Println("Commit:", commit.Hash, commit.Message)
		}
	}

	return nil
}
