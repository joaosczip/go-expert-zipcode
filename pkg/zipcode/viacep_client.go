package zipcode

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type ViaCepZipCodeClient struct {
	client *http.Client
}

func NewViaCepZipCodeClient(client *http.Client) *ViaCepZipCodeClient {
	return &ViaCepZipCodeClient{client}
}

func (v *ViaCepZipCodeClient) Fetch(context context.Context, zipCode string) (*City, error) {
	req, err := http.NewRequestWithContext(context, http.MethodGet, "https://viacep.com.br/ws/"+zipCode+"/json", nil)

	if err != nil {
		return nil, err
	}

	resp, err := v.client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrInvalidZipCode
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var viaCepResponse struct {
		Localidade string `json:"localidade"`
	}

	if err := json.Unmarshal(respBody, &viaCepResponse); err != nil {
		return nil, err
	}

	if viaCepResponse.Localidade == "" {
		return nil, ErrZipCodeNotFound
	}

	return &City{Name: viaCepResponse.Localidade}, nil
}
