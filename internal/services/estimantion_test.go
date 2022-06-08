package services

import (
	"fake-eta-task/internal/adapters"
	"testing"
)

type wheelyMock struct{}

func (w wheelyMock) GetCars(target adapters.Coordinates) ([]adapters.Car, error) {
	return []adapters.Car{
		{
			Id: 1,
			Coordinates: adapters.Coordinates{
				Lat: 55.7575429,
				Lng: 37.6135117,
			},
		},
		{
			Id: 2,
			Coordinates: adapters.Coordinates{
				Lat: 55.74837156167371,
				Lng: 37.61180107665421,
			},
		},
	}, nil
}

func (w wheelyMock) GetRoutePredictions(target adapters.Coordinates, source []adapters.Coordinates) ([]int, error) {
	return []int{1, 2}, nil
}

func TestNewEstimationService(t *testing.T) {
	wheely := wheelyMock{}
	estimationService := NewEstimationService(wheely)

	if estimationService == nil {
		t.Error("NewEstimationService() returned nil")
	}
}

func TestEstimate(t *testing.T) {
	wheely := wheelyMock{}
	estimationService := NewEstimationService(wheely)

	estimation, err := estimationService.Estimate(adapters.Coordinates{
		Lat: 0.0,
		Lng: 0.0,
	})

	if err != nil {
		t.Errorf("Estimate() returned error: %s", err.Error())
	}

	if estimation != 1 {
		t.Errorf("Estimate() returned %d instead of 1", estimation)
	}

}
