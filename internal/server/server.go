package server

import (
	"fake-eta-task/internal/services"
	"net/http"
)

type server struct {
	estimationService services.EstimationService
}

func NewServer(estimationService services.EstimationService) *server {
	return &server{
		estimationService: estimationService,
	}
}

func (s *server) InitRoutes() {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/getEstimation", s.handleGetEstimation)
}
