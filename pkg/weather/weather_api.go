package weather

import (
	"context"
	"encoding/json"
	"net/http"
)

type WeatherAPIClient struct {
	client *http.Client
}

type WeatherAPIErrorResponse struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

var (
	codeMissingLocation = 1003
)

func NewWeatherAPIClient(client *http.Client) *WeatherAPIClient {
	return &WeatherAPIClient{client}
}

func (w *WeatherAPIClient) Fetch(ctx context.Context, city string) (*Weather, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.weatherapi.com/v1/current.json?key=123&q="+city, nil)

	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		var errorResponse WeatherAPIErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return nil, err
		}

		if errorResponse.Error.Code == codeMissingLocation {
			return nil, ErrMissingLocation
		}

		return nil, ErrInvalidLocation
	}

	var weatherResponse struct {
		Current struct {
			TempC float64 `json:"temp_c"`
			TempF float64 `json:"temp_f"`
		} `json:"current"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return nil, err
	}

	weather := NewWeather(weatherResponse.Current.TempC, weatherResponse.Current.TempF)

	return weather, nil
}
