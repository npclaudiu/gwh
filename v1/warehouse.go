package gwh

import (
	"fmt"

	"github.com/npclaudiu/gwh/internal/control"
	"github.com/npclaudiu/gwh/internal/hostfs"
)

type Warehouse struct {
	controlDatabase *control.ControlDatabase
}

func Open(location string) (*Warehouse, error) {
	// Make sure we have a valid file layout.
	//
	warehouseLayout, err := hostfs.NewWarehouseLayout(location)

	if err != nil {
		return nil, fmt.Errorf("gwh: %w", err)
	}

	// Open control database.
	//
	controlDatabaseFile := warehouseLayout.GetPath(hostfs.ControlDatabaseFile)
	controlDatabase, err := control.OpenControlDatabase(controlDatabaseFile)

	if err != nil {
		return nil, fmt.Errorf("gwh: %w", err)
	}

	return &Warehouse{
		controlDatabase: controlDatabase,
	}, nil
}

func (w *Warehouse) LinkRepository(path string) error {
	return nil
}
