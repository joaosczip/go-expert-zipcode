package weather

import (
	"context"
	"encoding/json"
	"io"
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

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var weatherResponse struct {
		Current struct {
			TempC float64 `json:"temp_c"`
			TempF float64 `json:"temp_f"`
		} `json:"current"`
	}

	if err := json.Unmarshal(respBody, &weatherResponse); err != nil {
		return nil, err
	}

	weather := NewWeather(weatherResponse.Current.TempC, weatherResponse.Current.TempF)

	return weather, nil
}
