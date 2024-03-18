package zipcode

import (
	"context"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("should return a city when the zip code is found", func(t *testing.T) {
		mockedResponse := `
			{
				"cep":"01001-000",
				"logradouro":"Praça da Sééééééééeéééé",
				"complemento":"lado ímpar",
				"bairro":"Sé",
				"localidade":"São Paulo",
				"uf":"SP",
				"ibge":"3550308",
				"gia":"1004",
				"ddd":"11",
				"siafi":"7107"
			}
		`

		httpmock.RegisterResponder(
			http.MethodGet,
			"https://viacep.com.br/ws/01001000/json",
			httpmock.NewStringResponder(
				http.StatusOK,
				mockedResponse,
			),
		)

		client := NewViaCepZipCodeClient(http.DefaultClient)

		city, err := client.Fetch(context.TODO(), "01001000")

		assert.NoError(t, err)
		assert.Equal(t, city.Name, "São Paulo")
	})

	t.Run("should return an error when the zipcode was not found", func(t *testing.T) {
		mockedResponse := `
			{
				"erro": true
			}
		`

		httpmock.RegisterResponder(
			http.MethodGet,
			"https://viacep.com.br/ws/11111111/json",
			httpmock.NewStringResponder(
				http.StatusOK,
				mockedResponse,
			),
		)

		client := NewViaCepZipCodeClient(http.DefaultClient)

		city, err := client.Fetch(context.TODO(), "11111111")

		assert.NotNil(t, err)
		assert.Nil(t, city)
		assert.ErrorIs(t, err, ErrZipCodeNotFound)
	})

	t.Run("should return an error when an invlid zipcode was provided", func(t *testing.T) {
		httpmock.RegisterResponder(
			http.MethodGet,
			"https://viacep.com.br/ws/invalid-zip-code/json",
			httpmock.NewStringResponder(
				http.StatusBadRequest,
				"",
			),
		)

		client := NewViaCepZipCodeClient(http.DefaultClient)

		city, err := client.Fetch(context.TODO(), "invalid-zip-code")

		assert.NotNil(t, err)
		assert.Nil(t, city)
		assert.ErrorIs(t, err, ErrInvalidZipCode)
	})
}
