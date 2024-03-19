package weather

import (
	"context"
	"errors"
)

var ErrInvalidLocation = errors.New("invalid location")

type Weather struct {
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
	Kelvin     float64 `json:"kelvin"`
}

func NewWeather(celsius, fahrenheit float64) *Weather {
	return &Weather{
		Celsius:    celsius,
		Fahrenheit: fahrenheit,
		Kelvin:     celsius + 273,
	}
}

type WeatherClient interface {
	Fetch(ctx context.Context, city string) (*Weather, error)
}
