package adapters

import (
	"errors"
	"testing"
)

type mockCache struct {
}

func (m *mockCache) Get(key string, result interface{}) error {
	return errors.New("not cached")
}

func (m *mockCache) Set(key string, value interface{}) error {
	return nil
}

func TestWheely(t *testing.T) {
	cache := &mockCache{}
	wheely := NewWheely(cache)

	if wheely == nil {
		t.Error("NewWheely() returned nil")
	}
}

func TestGetCars(t *testing.T) {
	cache := &mockCache{}
	wheely := NewWheely(cache)

	cars, err := wheely.GetCars(Coordinates{
		Lat: 0.0,
		Lng: 0.0,
	}, 10)

	if err != nil {
		t.Errorf("GetCars() returned error: %s", err.Error())
	}

	if len(cars) != 10 {
		t.Errorf("GetCars() returned %d instead of 2", len(cars))
	}
}

func TestGetRoutePredictions(t *testing.T) {
	cache := &mockCache{}
	wheely := NewWheely(cache)

	route, err := wheely.GetRoutePredictions(Coordinates{
		Lat: 0.0,
		Lng: 0.0,
	}, []Coordinates{
		{
			Lat: 55.7575429,
			Lng: 37.6135117,
		},
		{
			Lat: 55.74837156167371,
			Lng: 37.61180107665421,
		},
	})

	if err != nil {
		t.Errorf("GetRoutePredictions() returned error: %s", err.Error())
	}

	if len(route) != 2 {
		t.Errorf("GetRoutePredictions() returned %d instead of 2", len(route))
	}
}
