package services

import (
	"errors"
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
	cars, err := adapters.Retry(3, 1, func() ([]adapters.Car, error) {
		return e.wheely.GetCars(target, 10)
	})

	if err != nil {
		return 0, err
	}

	if len(cars) == 0 {
		return 0, errors.New("No cars found")
	}

	sources := []adapters.Coordinates{}

	for _, car := range cars {
		sources = append(sources, adapters.Coordinates{
			Lat: car.Coordinates.Lat,
			Lng: car.Coordinates.Lng,
		})
	}

	predictions, err := adapters.Retry(3, 1, func() ([]int, error) {
		return e.wheely.GetRoutePredictions(target, sources)
	})

	if err != nil {
		return 0, err
	}

	min := predictions[0]
	for _, r := range predictions {
		if r < min {
			min = r
		}
	}

	return min, nil
}
