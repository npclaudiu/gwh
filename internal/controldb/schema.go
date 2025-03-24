package controldb

import (
	"embed"
)

//go:embed schema/*.sql
var embeddedSchemas embed.FS
