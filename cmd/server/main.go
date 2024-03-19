package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/joaosczip/zipcode_temp/internal/routes"
)

func main() {
	r := chi.NewRouter()

	r.Get("/temperature", routes.FetchTemperature)

	http.ListenAndServe(":8080", r)
}
