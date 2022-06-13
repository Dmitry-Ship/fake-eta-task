package services

import (
	"errors"
	"fake-eta-task/internal/common"
	carsClient "fake-eta-task/internal/generated/cars/client/operations"
	carsModels "fake-eta-task/internal/generated/cars/models"
	predictionClient "fake-eta-task/internal/generated/prediction/client/operations"
	"fake-eta-task/internal/generated/prediction/models"

	"math/rand"

	"fake-eta-task/internal/infra"
	"fmt"
)

type EstimationService interface {
	Estimate(lat, lng float64) (int64, error)
}

type estimationService struct {
	carsService       carsClient.ClientService
	predictionService predictionClient.ClientService
	cache             infra.Cache
}

func NewEstimationService(carsService carsClient.ClientService, predictionService predictionClient.ClientService, cache infra.Cache) *estimationService {
	return &estimationService{
		carsService:       carsService,
		predictionService: predictionService,
		cache:             cache,
	}
}

func (s estimationService) Estimate(lat, lng float64) (int64, error) {
	var minimalTime int64 = 0
	cacheKey := "time_prediction_" + fmt.Sprintf("%f", lat) + "_" + fmt.Sprintf("%f", lng)
	err := s.cache.Get(cacheKey, &minimalTime)

	if err == nil {
		return minimalTime, nil
	}

	carsResponse, err := common.Retry(2, 1, func() (*carsClient.GetCarsOK, error) {
		params := carsClient.NewGetCarsParams()
		params.SetLat(lat)
		params.SetLng(lng)
		params.SetLimit(10)

		return s.carsService.GetCars(params)
	})

	// ⚠️🤷🏻‍♂️ Normally it would just return error, but since it is a fake service, we need to return some fake cars
	if err != nil {
		carsResponse = &carsClient.GetCarsOK{
			Payload: s.generateFakeCars(lat, lng, 10),
		}
	}

	if len(carsResponse.Payload) == 0 {
		return 0, errors.New("No cars found")
	}

	sources := []models.Position{}

	for _, car := range carsResponse.Payload {
		sources = append(sources, models.Position{
			Lat: car.Lat,
			Lng: car.Lng,
		})
	}

	predictionResponse, err := common.Retry(2, 1, func() (*predictionClient.PredictOK, error) {
		params := predictionClient.NewPredictParams()
		predictBody := predictionClient.PredictBody{
			Source: sources,
			Target: models.Position{
				Lat: lat,
				Lng: lng,
			},
		}
		params.SetPositionList(predictBody)

		return s.predictionService.Predict(params)
	})

	// ⚠️🤷🏻‍♂️ Normally it would just return error, but since it is a fake service, we need to return some fake time
	if err != nil {
		predictionResponse = &predictionClient.PredictOK{
			Payload: s.generateFakePredictions(10),
		}
	}

	min := predictionResponse.Payload[0]
	for _, r := range predictionResponse.Payload {
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

func (s estimationService) generateFakeCars(lat float64, lng float64, numberOfCars int) []carsModels.Car {
	cars := []carsModels.Car{}
	for i := 0; i < numberOfCars; i++ {
		car := carsModels.Car{
			ID:  int64(i),
			Lat: lat + rand.Float64()/100,
			Lng: lng + rand.Float64()/100,
		}
		cars = append(cars, car)
	}

	return cars
}

func (s estimationService) generateFakePredictions(numberOfSources int) []int64 {
	predictions := []int64{}
	for i := 0; i < numberOfSources; i++ {
		// let's assume travel time is never longer than 30 minutes
		predictions = append(predictions, int64(1+rand.Intn(30)))
	}

	return predictions
}
