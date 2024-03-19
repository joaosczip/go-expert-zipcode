package weather

import (
	"context"
	"encoding/json"
	"net/http"
)

type WeatherAPIClient struct {
	client *http.Client
}

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

	if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrInvalidLocation
	}

	defer resp.Body.Close()

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
