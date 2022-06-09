package main

import (
	"fake-eta-task/internal/adapters"
	"fake-eta-task/internal/server"
	"fake-eta-task/internal/services"
	"log"
	"net/http"
	"os"
)

func main() {
	wheely := adapters.NewWheely()
	estimationService := services.NewEstimationService(wheely)
	server := server.NewServer(estimationService)

	server.InitRoutes()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
