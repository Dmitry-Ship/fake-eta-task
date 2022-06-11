package main

import (
	"context"
	"fake-eta-task/internal/adapters"
	"fake-eta-task/internal/infra"
	"fake-eta-task/internal/server"
	"fake-eta-task/internal/services"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	cache := infra.NewCache(ctx)
	wheely := adapters.NewWheely()
	estimationService := services.NewEstimationService(wheely, cache)
	server := server.NewServer(estimationService)

	server.InitRoutes()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
