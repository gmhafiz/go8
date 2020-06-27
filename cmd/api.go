package main

import (
	"fmt"
	"log"
	"net/http"

	"eight/app"
	"eight/config"
)

func main() {
	c := config.AppConfig()
	newApp := app.NewApp(c)
	router := newApp.Router()

	address := fmt.Sprintf(":%d", newApp.GetConfig().Server.Port)

	server := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  newApp.GetConfig().Server.TimeoutRead,
		WriteTimeout: newApp.GetConfig().Server.TimeoutWrite,
		IdleTimeout:  newApp.GetConfig().Server.TimeoutIdle,
	}

	log.Printf("starting server %v", address)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("fails to start server")
	}
}


