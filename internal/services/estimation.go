package services

import (
	"errors"
	"fake-eta-task/internal/adapters"
	"fake-eta-task/internal/infra"
	"fmt"
)

type EstimationService interface {
	Estimate(adapters.Coordinates) (int, error)
}

type estimationService struct {
	wheely adapters.Wheely
	cache  infra.Cache
}

func NewEstimationService(wheely adapters.Wheely, cache infra.Cache) *estimationService {
	return &estimationService{
		wheely: wheely,
		cache:  cache,
	}
}

func (s estimationService) Estimate(target adapters.Coordinates) (int, error) {
	minimalTime := 0
	cacheKey := "time_prediction_" + fmt.Sprintf("%f", target.Lat) + "_" + fmt.Sprintf("%f", target.Lng)
	err := s.cache.Get(cacheKey, &minimalTime)

	if err == nil {
		return minimalTime, nil
	}

	cars, err := adapters.Retry(3, 1, func() ([]adapters.Car, error) {
		return s.wheely.GetCars(target, 10)
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
		return s.wheely.GetRoutePredictions(target, sources)
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

	err = s.cache.Set(cacheKey, 10, min)

	if err != nil {
		return 0, err
	}

	return min, nil
}
