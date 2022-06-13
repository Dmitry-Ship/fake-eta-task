package server

import (
	"fmt"
	"net/http"
	"strconv"
)

func (s *server) handleGetEstimation(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lat := r.URL.Query().Get("lat")
	if lat == "" {
		http.Error(w, "Target coordinates are required", http.StatusBadRequest)
		return
	}

	lng := r.URL.Query().Get("lng")
	if lng == "" {
		http.Error(w, "Target coordinates are required", http.StatusBadRequest)
		return
	}

	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		http.Error(w, "Invalid target coordinates", http.StatusBadRequest)
		return
	}

	lngFloat, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		http.Error(w, "Invalid target coordinates", http.StatusBadRequest)
		return
	}

	estimation, err := s.estimationService.Estimate(latFloat, lngFloat)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(fmt.Sprintf("%d", estimation)))

	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %s", err.Error()), http.StatusInternalServerError)
		return
	}

}

func (s *server) handleIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "internal/static/index.html")
}
