package main

import (
	"log"

	"github.com/gmhafiz/go8/internal/server"
)

const Version = "v0.7.0"

// @title Go8
// @version 0.7.0
// @description Go + Postgres + Chi Router + sqlx + Unit Testing starter kit for API development.

// @contact.name Hafiz Shafruddin
// @contact.url https://github.com/gmhafiz/go8
// @contact.email gmhafiz@gmail.com

// @host localhost:3080
// @BasePath /
func main() {
	s := server.New(Version)
	s.Init()

	if err := s.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
