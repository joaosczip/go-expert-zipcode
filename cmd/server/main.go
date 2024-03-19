package main

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/joaosczip/zipcode_temp/internal/handlers"
	"github.com/joaosczip/zipcode_temp/pkg/weather"
	"github.com/joaosczip/zipcode_temp/pkg/zipcode"
)

type httpResponse struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func main() {
	http.HandleFunc("/temperature", func(w http.ResponseWriter, r *http.Request) {
		httpClient := http.DefaultClient

		zipcodeClient := zipcode.NewViaCepZipCodeClient(httpClient)
		weatherClient := weather.NewWeatherAPIClient(httpClient, "905a326931cb4138a99225100241803")
		handler := handlers.NewTemperatureHandler(zipcodeClient, weatherClient)

		zcode := r.URL.Query().Get("zipcode")

		output, err := handler.FetchTemperature(zcode)

		var response httpResponse

		if err != nil {
			if errors.Is(err, zipcode.ErrInvalidZipCode) {
				response = httpResponse{
					StatusCode: http.StatusUnprocessableEntity,
					Message:    err.Error(),
				}
			} else if errors.Is(err, zipcode.ErrZipCodeNotFound) {
				response = httpResponse{
					StatusCode: http.StatusNotFound,
					Message:    err.Error(),
				}
			} else {
				response = httpResponse{
					StatusCode: http.StatusInternalServerError,
					Message:    err.Error(),
				}
			}
			w.WriteHeader(response.StatusCode)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(output); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", nil)
}
