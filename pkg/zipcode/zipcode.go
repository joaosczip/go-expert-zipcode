package zipcode

import (
	"context"
	"errors"
)

var (
	ErrZipCodeNotFound = errors.New("can not find zipcode")
	ErrInvalidZipCode  = errors.New("invalid zipcode")
)

type City struct {
	Name string `json:"name"`
}

type ZipCodeClient interface {
	Fetch(ctx context.Context, zipCode string) (*City, error)
}
