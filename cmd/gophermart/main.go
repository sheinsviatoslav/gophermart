package main

import (
	"github.com/sheinsviatoslav/gophermart/internal/config"
	"github.com/sheinsviatoslav/gophermart/internal/routes"
	"log"
	"net/http"
)

func main() {
	config.Init()

	log.Println("listen on", *config.RunAddress)
	log.Fatal(http.ListenAndServe(*config.RunAddress, routes.MainRouter()))
}
