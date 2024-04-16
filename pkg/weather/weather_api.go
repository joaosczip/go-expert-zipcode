package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type WeatherAPIClient struct {
	client *http.Client
	apiKey string
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

func NewWeatherAPIClient(client *http.Client, apiKey string) *WeatherAPIClient {
	return &WeatherAPIClient{client, apiKey}
}

func (w *WeatherAPIClient) Fetch(ctx context.Context, city string) (*Weather, error) {
	reqUrl := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", w.apiKey, strings.ReplaceAll(city, " ", "+"))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqUrl, nil)

	if err != nil {
		return nil, err
	}

	resp, err := w.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError {
		if resp.StatusCode == http.StatusForbidden {
			return nil, ErrInvalidAPIKey
		}

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
