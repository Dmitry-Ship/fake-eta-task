package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
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

func (s *wheelyService) GetCars(target Coordinates, numberOfCars int) ([]Car, error) {
	resp, err := http.Get("https://dev-api.wheely.com/fake-eta/cars?lat=" + fmt.Sprintf("%f", target.Lat) + "&lng=" + fmt.Sprintf("%f", target.Lng) + "&limit=" + fmt.Sprintf("%d", numberOfCars))

	if err != nil {
		return []Car{}, err
	}

	defer resp.Body.Close()

	// ⚠️🤷🏻‍♂️ Normally it would just return error, but since it is a fake service, we need to return some fake cars
	if resp.StatusCode != http.StatusOK {
		cars := []Car{}
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

		return cars, nil
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []Car{}, err
	}

	cars := []Car{}

	if err := json.Unmarshal(body, &cars); err != nil {
		return []Car{}, err
	}

	return cars, nil
}

func (s *wheelyService) GetRoutePredictions(target Coordinates, source []Coordinates) ([]int, error) {
	req := struct {
		Target Coordinates   `json:"target"`
		Source []Coordinates `json:"source"`
	}{
		Target: target,
		Source: source,
	}

	json_data, err := json.Marshal(req)

	if err != nil {
		return []int{}, err
	}

	resp, err := http.Post("https://dev-api.wheely.com/fake-eta/predict", "application/json", bytes.NewBuffer(json_data))

	if err != nil {
		return []int{}, err
	}

	defer resp.Body.Close()

	// ⚠️🤷🏻‍♂️ Normally it would just return error, but since it is a fake service, we need to return some fake time
	if resp.StatusCode != http.StatusOK {
		travelTimeFromSourceToTarget := []int{}

		for range source {
			// let's assume travel time is never longer than 30 minutes
			travelTimeFromSourceToTarget = append(travelTimeFromSourceToTarget, 1+rand.Intn(30))
		}

		return travelTimeFromSourceToTarget, nil
	}

	var res []int

	err = json.NewDecoder(resp.Body).Decode(&res)

	if err != nil {
		return []int{}, err
	}

	return res, nil
}
