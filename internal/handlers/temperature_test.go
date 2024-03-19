package handlers

import (
	"context"
	"testing"

	"github.com/joaosczip/zipcode_temp/pkg/weather"
	"github.com/joaosczip/zipcode_temp/pkg/zipcode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ZipCodeClientMock struct {
	mock.Mock
}

func (m *ZipCodeClientMock) Fetch(ctx context.Context, zipCode string) (*zipcode.City, error) {
	args := m.Called(zipCode)
	return args.Get(0).(*zipcode.City), args.Error(1)
}

type WeatherClientMock struct {
	mock.Mock
}

func (m *WeatherClientMock) Fetch(ctx context.Context, city string) (*weather.Weather, error) {
	args := m.Called(city)
	return args.Get(0).(*weather.Weather), args.Error(1)
}

func TestTemperatureHandler_FetchTemperature(t *testing.T) {
	t.Run("should return the temperature when the zip code is found", func(t *testing.T) {
		zipCodeClientMock := &ZipCodeClientMock{}
		weatherClientMock := &WeatherClientMock{}

		zipCodeClientMock.On("Fetch", "01001000").Return(&zipcode.City{
			Name: "São Paulo",
		}, nil)

		weatherClientMock.On("Fetch", "São Paulo").Return(&weather.Weather{
			Celsius:    22.0,
			Fahrenheit: 71.6,
			Kelvin:     295.0,
		}, nil)

		handler := NewTemperatureHandler(zipCodeClientMock, weatherClientMock)

		output, err := handler.FetchTemperature("01001000")

		assert.Nil(t, err)
		assert.Equal(t, output.Celsius, 22.0)
		assert.Equal(t, output.Fahrenheit, 71.6)
		assert.Equal(t, output.Kelvin, 295.0)

		zipCodeClientMock.AssertExpectations(t)
		weatherClientMock.AssertExpectations(t)
	})

	t.Run("should return an error when the city is not found", func(t *testing.T) {
		zipCodeClientMock := &ZipCodeClientMock{}
		weatherClientMock := &WeatherClientMock{}

		zipCodeClientMock.On("Fetch", "11111111").Return(&zipcode.City{}, zipcode.ErrZipCodeNotFound)

		handler := NewTemperatureHandler(zipCodeClientMock, weatherClientMock)

		_, err := handler.FetchTemperature("11111111")

		assert.Error(t, err)
		assert.ErrorIs(t, err, zipcode.ErrZipCodeNotFound)

		zipCodeClientMock.AssertExpectations(t)
		weatherClientMock.AssertNotCalled(t, "Fetch")
	})
}
