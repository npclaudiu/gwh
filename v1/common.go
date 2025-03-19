package gwh

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	WarehouseDirName          = ".gwh"
	ManifestFileName          = "manifest.yaml"
	AnalyticsDatabaseFileName = "analytics.db"
)

var (
	ErrInvalidPrefix          = errors.New("invalid prefix")
	ErrWarehouseDirIsInvalid  = errors.New("warehouse directory is invalid")
	ErrWarehouseDirExists     = errors.New("warehouse directory already exists")
	ErrWarehouseDirInitFailed = errors.New("failed to initialize warehouse directory")
)

func validatePrefix(prefix string) (string, error) {
	if !filepath.IsAbs(prefix) {
		return "", ErrInvalidPrefix
	}

	if stat, err := os.Stat(prefix); err != nil || !stat.IsDir() {
		return "", ErrInvalidPrefix
	}

	return prefix, nil
}
