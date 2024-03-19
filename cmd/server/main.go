package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/joaosczip/zipcode_temp/configs"
	"github.com/joaosczip/zipcode_temp/internal/routes"
)

func main() {
	_, err := configs.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()

	r.Get("/temperature", routes.FetchTemperature)

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	} else {
		println("Server running on port 8080")
	}
}
