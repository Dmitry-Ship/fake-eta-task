package services

import (
	"fake-eta-task/internal/adapters"
)

type EstimationService interface {
	Estimate(adapters.Coordinates) (int, error)
}

type estimationService struct {
	wheely adapters.Wheely
}

func NewEstimationService(wheely adapters.Wheely) *estimationService {
	return &estimationService{
		wheely: wheely,
	}
}

func (e estimationService) Estimate(target adapters.Coordinates) (int, error) {
	cars, err := e.wheely.GetCars(target)

	if err != nil {
		return 0, err
	}

	source := []adapters.Coordinates{}

	for _, car := range cars {
		source = append(source, adapters.Coordinates{
			Lat: car.Coordinates.Lat,
			Lng: car.Coordinates.Lng,
		})
	}

	route, err := e.wheely.GetRoutePredictions(target, source)

	if err != nil {
		return 0, err
	}

	min := route[0]
	for _, r := range route {
		if r < min {
			min = r
		}
	}

	return min, nil
}
