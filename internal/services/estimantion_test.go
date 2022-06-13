package services

import (
	"errors"
	carsClient "fake-eta-task/internal/generated/cars/client/operations"
	predictionClient "fake-eta-task/internal/generated/prediction/client/operations"

	"fake-eta-task/internal/generated/cars/models"
	"testing"
	"time"

	"github.com/go-openapi/runtime"
)

type mockCache struct {
}

func (m *mockCache) Get(key string, result interface{}) error {
	return errors.New("not cached")
}

func (m *mockCache) Set(key string, lifeTime time.Duration, value interface{}) error {
	return nil
}

type carsClientMock struct{}

func (c carsClientMock) GetCars(params *carsClient.GetCarsParams, opts ...carsClient.ClientOption) (*carsClient.GetCarsOK, error) {
	return &carsClient.GetCarsOK{
		Payload: []models.Car{
			{
				ID:  1,
				Lat: 55.7575429,
				Lng: 37.6135117,
			},
			{
				ID:  2,
				Lat: 55.74837156167371,
				Lng: 37.61180107665421,
			},
		}}, nil
}

func (c carsClientMock) Health(params *carsClient.HealthParams, opts ...carsClient.ClientOption) (*carsClient.HealthOK, error) {
	return &carsClient.HealthOK{}, nil
}

func (c carsClientMock) SetTransport(transport runtime.ClientTransport) {
}

type predictionClientMock struct{}

func (p predictionClientMock) Predict(params *predictionClient.PredictParams, opts ...predictionClient.ClientOption) (*predictionClient.PredictOK, error) {
	return &predictionClient.PredictOK{
		Payload: []int64{1, 2, 3},
	}, nil
}

func (p predictionClientMock) Health(params *predictionClient.HealthParams, opts ...predictionClient.ClientOption) (*predictionClient.HealthOK, error) {
	return &predictionClient.HealthOK{}, nil
}

func (p predictionClientMock) SetTransport(transport runtime.ClientTransport) {
}

func TestNewEstimationService(t *testing.T) {
	carsClient := carsClientMock{}
	predictionClient := predictionClientMock{}
	cache := &mockCache{}
	estimationService := NewEstimationService(carsClient, predictionClient, cache)

	if estimationService == nil {
		t.Error("NewEstimationService() returned nil")
	}
}

func TestEstimate(t *testing.T) {
	carsClient := carsClientMock{}
	predictionClient := predictionClientMock{}
	cache := &mockCache{}
	estimationService := NewEstimationService(carsClient, predictionClient, cache)

	estimation, err := estimationService.Estimate(0.0, 0.0)

	if err != nil {
		t.Errorf("Estimate() returned error: %s", err.Error())
	}

	if estimation != 1 {
		t.Errorf("Estimate() returned %d instead of 1", estimation)
	}
}
