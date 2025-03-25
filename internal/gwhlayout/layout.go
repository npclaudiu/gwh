package gwhlayout

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const (
	WarehouseDirectoryName  = ".gwh"
	ControlDatabaseFileName = "control.db"
)

type WarehouseLayoutPath int

const (
	WorkingDirectory WarehouseLayoutPath = iota
	WarehouseDirectory
	ControlDatabaseFile
)

var (
	nameRegexp = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`)
)

type WarehouseLayout struct {
	workingDirectoryPath    string
	warehouseDirectoryPath  string
	controlDatabaseFilePath string
}

func New(location string) (*WarehouseLayout, error) {
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
		workingDirectoryPath:    location,
		warehouseDirectoryPath:  whDir,
		controlDatabaseFilePath: path.Join(whDir, ControlDatabaseFileName),
	}

	return wl, nil
}

func (wl *WarehouseLayout) GetKnownPath(p WarehouseLayoutPath) string {
	switch p {
	case WorkingDirectory:
		return wl.workingDirectoryPath
	case WarehouseDirectory:
		return wl.warehouseDirectoryPath
	case ControlDatabaseFile:
		return wl.controlDatabaseFilePath
	default:
		return ""
	}
}

func (wl *WarehouseLayout) ResolvePath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}

	return path.Join(wl.workingDirectoryPath, p)
}

func (wl *WarehouseLayout) IsNameValid(name string) bool {
	return nameRegexp.MatchString(name)
}
