package gwh

import (
	"errors"
)

const (
	WarehouseDirName = ".gwh"
	ManifestFileName = "manifest.yaml"
	DatabaseFileName = "main.db"
)

var (
	ErrWorkDirInvalid         = errors.New("invalid working directory")
	ErrWarehouseDirExists     = errors.New("warehouse directory already exists")
	ErrWarehouseDirInitFailed = errors.New("failed to initialize warehouse directory")
)
