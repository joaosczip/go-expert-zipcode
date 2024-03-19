package handlers

import (
	"context"
	"time"

	"github.com/joaosczip/zipcode_temp/pkg/weather"
	"github.com/joaosczip/zipcode_temp/pkg/zipcode"
)

type TemperatureHandler struct {
	zipCodeClient zipcode.ZipCodeClient
	weatherClient weather.WeatherClient
}

func NewTemperatureHandler(zipCodeClient zipcode.ZipCodeClient, weatherClient weather.WeatherClient) *TemperatureHandler {
	return &TemperatureHandler{zipCodeClient: zipCodeClient, weatherClient: weatherClient}
}

type FetchTemperatureOutput struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

func (h *TemperatureHandler) FetchTemperature(zipCode string) (FetchTemperatureOutput, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	city, err := h.zipCodeClient.Fetch(ctx, zipCode)

	if err != nil {
		return FetchTemperatureOutput{}, err
	}

	weather, err := h.weatherClient.Fetch(ctx, city.Name)

	if err != nil {
		return FetchTemperatureOutput{}, err
	}

	return FetchTemperatureOutput{
		Celsius:    weather.Celsius,
		Fahrenheit: weather.Fahrenheit,
		Kelvin:     weather.Kelvin,
	}, nil
}
