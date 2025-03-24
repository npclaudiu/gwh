package hostfs

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
)

const (
	WarehouseDirectoryName  = ".gwh"
	ControlDatabaseFileName = "control.db"
)

type WarehouseLayoutPath int

const (
	WarehouseDirectory WarehouseLayoutPath = iota
	ControlDatabaseFile
)

type WarehouseLayout struct {
	warehouseDirectoryPath  string
	controlDatabaseFilePath string
}

func NewWarehouseLayout(location string) (*WarehouseLayout, error) {
	// Check location.
	//
	cwd, err := os.Getwd()

	if err != nil {
		return nil, fmt.Errorf("gwh: failed to load working directory: %w", err)
	}

	if !filepath.IsAbs(location) {
		location = filepath.Join(cwd, location)
	}

	if stat, err := os.Stat(location); err != nil || !stat.IsDir() {
		return nil, fmt.Errorf("gwh: failed to open warehouse: %w", err)
	}

	// Check warehouse directory, creating it if it doesn't exist.
	//
	whDir := path.Join(location, WarehouseDirectoryName)
	_, err = os.Stat(whDir)
	whDirExists := true

	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			whDirExists = false
		} else {
			return nil, fmt.Errorf("gwh: failed to stat warehouse path: %w", err)
		}
	}

	if !whDirExists {
		err := os.Mkdir(whDir, 0755)

		if err != nil {
			return nil, fmt.Errorf("gwh: failed to create warehouse directory: %w", err)
		}
	}

	wl := &WarehouseLayout{
		warehouseDirectoryPath:  whDir,
		controlDatabaseFilePath: path.Join(whDir, ControlDatabaseFileName),
	}

	return wl, nil
}

func (wl *WarehouseLayout) GetPath(p WarehouseLayoutPath) string {
	switch p {
	case WarehouseDirectory:
		return wl.warehouseDirectoryPath
	case ControlDatabaseFile:
		return wl.controlDatabaseFilePath
	default:
		return ""
	}
}
