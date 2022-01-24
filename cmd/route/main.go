package main

import (
	"fmt"

	"github.com/gmhafiz/go8/internal/server"
)

// Version is injected using ldflags during build time
var Version = "v0.1.0"

func main() {
	s := server.New()
	s.InitDomains()
	fmt.Printf("Registered Routes:\n\n")
	s.PrintAllRegisteredRoutes()
}
