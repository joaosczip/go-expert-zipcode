package weather

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestWeatherAPIClient(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("should return the weather when the city is found", func(t *testing.T) {
		mockedResponse := `
			{
				"current": {
					"temp_c": 22.0,
					"temp_f": 71.6
				}
			}
		`

		httpmock.RegisterResponder(
			http.MethodGet,
			"https://api.weatherapi.com/v1/current.json?key=123&q=Curitiba",
			httpmock.NewStringResponder(
				http.StatusOK,
				mockedResponse,
			),
		)

		weatherClient := NewWeatherAPIClient(http.DefaultClient, "123")

		weather, err := weatherClient.Fetch(context.TODO(), "Curitiba")

		assert.NoError(t, err)
		assert.Equal(t, weather.Celsius, 22.0)
		assert.Equal(t, weather.Fahrenheit, 71.6)
		assert.Equal(t, weather.Kelvin, 295.0)
	})

	t.Run("should return an error when the location was not found", func(t *testing.T) {
		mockedResponse := `
			{
				"error": {
					"code": 1006,
					"message": "No matching location found."
				}
			}
		`

		httpmock.RegisterResponder(
			http.MethodGet,
			"https://api.weatherapi.com/v1/current.json?key=123&q=invalid",
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				mockedResponse,
			),
		)

		weatherClient := NewWeatherAPIClient(http.DefaultClient, "123")

		weather, err := weatherClient.Fetch(context.TODO(), "invalid")

		assert.Error(t, err)
		assert.Nil(t, weather)
		assert.ErrorIs(t, err, ErrInvalidLocation)
	})

	t.Run("should return an error when the city is blank", func(t *testing.T) {
		mockedResponse := `
			{
				"error": {
					"code": 1003,
					"message": "Parameter q is missing."
				}
			}
		`

		httpmock.RegisterResponder(
			http.MethodGet,
			"https://api.weatherapi.com/v1/current.json?key=123&q=",
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				mockedResponse,
			),
		)

		weatherClient := NewWeatherAPIClient(http.DefaultClient, "123")

		weather, err := weatherClient.Fetch(context.TODO(), "")

		assert.Error(t, err)
		assert.Nil(t, weather)
		assert.ErrorIs(t, err, ErrMissingLocation)
	})

	t.Run("should return an error when the api key is blank", func(t *testing.T) {
		httpmock.RegisterResponder(
			http.MethodGet,
			"https://api.weatherapi.com/v1/current.json?key=&q=Curitiba",
			httpmock.NewStringResponder(
				http.StatusForbidden,
				"",
			),
		)

		weatherClient := NewWeatherAPIClient(http.DefaultClient, "")

		weather, err := weatherClient.Fetch(context.TODO(), "Curitiba")

		assert.Error(t, err)
		assert.Nil(t, weather)
		assert.ErrorIs(t, err, ErrInvalidAPIKey)
	})
}
