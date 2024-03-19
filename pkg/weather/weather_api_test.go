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

		weatherClient := NewWeatherAPIClient(http.DefaultClient)

		weather, err := weatherClient.Fetch(context.TODO(), "Curitiba")

		assert.NoError(t, err)
		assert.Equal(t, weather.Celsius, 22.0)
		assert.Equal(t, weather.Fahrenheit, 71.6)
		assert.Equal(t, weather.Kelvin, 295.0)
	})
}
