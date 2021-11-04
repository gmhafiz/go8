package main

import (
	"log"

	"github.com/gmhafiz/go8/cmd/extmigrate/migrate"
)

// Version is injected using ldflags during build time
var Version string

func main() {
	log.Printf("Version: %s\n", Version)
	migrate.Start()
}
