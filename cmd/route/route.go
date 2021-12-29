package main

import (
	"fmt"
	"log"

	"github.com/gmhafiz/go8/internal/server"
)

// Version is injected using ldflags during build time
var Version = "v0.1.0"

func main() {
	log.Printf("Starting API version: %s\n", Version)
	s := server.New()
	s.InitDomains()
	fmt.Println("Registered Routes:")
	s.PrintAllRegisteredRoutes()
}
