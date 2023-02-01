package main

import (
	"log"

	"github.com/gmhafiz/go8/database"
	db "github.com/gmhafiz/go8/third_party/database"

	"github.com/gmhafiz/go8/config"
)

// Version is injected using ldflags during build time
var Version string

func main() {
	log.Printf("Version: %s\n", Version)

	cfg := config.New()
	store := db.NewSqlx(cfg.Database)
	migrator := database.Migrator(database.WithDB(store.DB))

	// todo: accept cli flag for other operations
	migrator.Up()
}
