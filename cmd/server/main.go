package main

import (
	"context"
	carsClient "fake-eta-task/internal/generated/cars/client"
	predictionClient "fake-eta-task/internal/generated/prediction/client"

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

	cars := carsClient.Default
	prediction := predictionClient.Default
	estimationService := services.NewEstimationService(cars.Operations, prediction.Operations, cache)
	server := server.NewServer(estimationService)

	server.InitRoutes()

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
