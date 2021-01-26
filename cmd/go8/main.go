package main

import (
	"log"

	"github.com/gmhafiz/go8/internal/server"
)

const Version = "v0.5.0"

func main() {
	s := server.New(Version)
	s.Init()

	if err := s.Run(); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
