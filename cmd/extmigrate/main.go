package main

import (
	"log"

	"github.com/gmhafiz/go8/database"
)

// Version is injected using ldflags during build time
var Version string

func main() {
	log.Printf("Version: %s\n", Version)

	migrator := database.Migrator()

	// todo: accept cli flag for other operations
	migrator.Up()
}
