package main

import (
	"fmt"

	"github.com/gmhafiz/go8/internal/server"
)

func main() {
	s := server.New()
	s.InitDomains()
	fmt.Print("Registered Routes:\n\n")
	s.PrintAllRegisteredRoutes()
}
