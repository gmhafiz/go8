package main

import (
	"log"
	"net/http"

	"eight/app"
)

func main() {
	srv := app.NewApp()
	r := srv.RegisterRoutes()

	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatalln("fails to start server")
	}
}


