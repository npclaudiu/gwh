package gwh

import (
	"fmt"

	"github.com/npclaudiu/gwh/internal/controldb"
	"github.com/npclaudiu/gwh/internal/hostfs"
)

type Warehouse struct {
	warehouseLayout *hostfs.WarehouseLayout
	controlDatabase *controldb.ControlDatabase
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
	controlDatabase, err := controldb.Open(controlDatabaseFile)

	if err != nil {
		return nil, fmt.Errorf("gwh: failed to open control database: %w", err)
	}

	return &Warehouse{
		warehouseLayout: warehouseLayout,
		controlDatabase: controlDatabase,
	}, nil
}

func (w *Warehouse) Close() {
	w.controlDatabase.Close()
}
