package main

import (
	"log"

	"github.com/gmhafiz/go8/internal/server"
)

// Version is injected using ldflags during build time
var Version = "v0.13.0"

// @title Go8
// @version 0.12.0
// @description Go + Postgres + Chi router + sqlx + ent + Unit Testing starter kit for API development.
// @contact.name Hafiz Shafruddin
// @contact.url https://github.com/gmhafiz/go8
// @contact.email gmhafiz@gmail.com
// @host localhost:3080
// @BasePath /
func main() {
	log.Printf("Starting API version: %s\n", Version)
	s := server.New()
	s.Init(Version)
	s.Run()
}
