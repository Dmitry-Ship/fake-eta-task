package adapters

import (
	"math/rand"
	"time"
)

type Wheely interface {
	GetCars(target Coordinates, numberOfCars int) ([]Car, error)
	GetRoutePredictions(target Coordinates, source []Coordinates) ([]int, error)
}

type wheelyService struct {
}

func NewWheely() *wheelyService {
	return &wheelyService{}
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Car struct {
	Coordinates
	Id int `json:"id"`
}

// 🤔 This what wheely API would return if I understood the task correctly.
func (c wheelyService) GetCars(target Coordinates, numberOfCars int) ([]Car, error) {
	cars := []Car{}

	// generate random cars around the target
	for i := 0; i < numberOfCars; i++ {
		car := Car{
			Id: i,
			Coordinates: Coordinates{
				Lat: target.Lat + rand.Float64()/100,
				Lng: target.Lng + rand.Float64()/100,
			},
		}
		cars = append(cars, car)
	}

	// 🤷🏻‍♂️ imitate latency for realism
	time.Sleep(time.Millisecond * 100)

	return cars, nil
}

// 🤔 This what wheely API would return if I understood the task correctly.
func (c wheelyService) GetRoutePredictions(target Coordinates, source []Coordinates) ([]int, error) {
	travelTimeFromSourceToTarget := []int{}

	for range source {
		// let's assume travel time is never longer than 30 minutes
		travelTimeFromSourceToTarget = append(travelTimeFromSourceToTarget, 1+rand.Intn(30))
	}

	// 🤷🏻‍♂️ imitate latency for realism
	time.Sleep(time.Millisecond * 100)

	return travelTimeFromSourceToTarget, nil
}
