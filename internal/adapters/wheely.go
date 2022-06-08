package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Wheely interface {
	GetCars(target Coordinates) ([]Car, error)
	GetRoutePredictions(target Coordinates, source []Coordinates) ([]int, error)
}

type wheelyService struct {
}

func NewWheely() wheelyService {
	return wheelyService{}
}

type Coordinates struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Car struct {
	Coordinates
	Id int `json:"id"`
}

func (c wheelyService) GetCars(target Coordinates) ([]Car, error) {
	resp, err := http.Get("https://dev-api.wheely.com/fake-eta/cars?lat=" + fmt.Sprintf("%f", target.Lat) + "&lng=" + fmt.Sprintf("%f", target.Lng) + "&limit=3")

	if err != nil {
		return []Car{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return []Car{}, err
	}

	cars := []Car{}

	if err := json.Unmarshal(body, &cars); err != nil {
		log.Println("err", err)

		return []Car{}, err
	}

	return cars, nil
}

func (c wheelyService) GetRoutePredictions(target Coordinates, source []Coordinates) ([]int, error) {
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

	var res []int

	json.NewDecoder(resp.Body).Decode(&res)

	return res, nil
}
