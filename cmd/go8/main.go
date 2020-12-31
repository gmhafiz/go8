package main

import (
	"log"

	"github.com/gmhafiz/go8/configs"
	"github.com/gmhafiz/go8/internal/server"
)

const Version = "v0.3.0"

func main() {
	cfg := configs.New()

	app := server.NewApp(cfg)

	if err := app.Run(cfg, Version); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
