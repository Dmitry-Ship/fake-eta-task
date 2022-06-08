package server

import (
	"fake-eta-task/internal/adapters"
	"fmt"
	"net/http"
)

func (s *server) handleGetEstimation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	estimation, err := s.estimationService.Estimate(adapters.Coordinates{
		Lat: 0.0,
		Lng: 0.0,
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	w.Write([]byte(fmt.Sprintf("%d", estimation)))
}
