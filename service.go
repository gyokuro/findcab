package findcab

import (
	"errors"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type Cab struct {
	Id        string
	Latitude  float64
	Longitude float64
}

type CabService interface {
	Read(id string) (Cab, error)
	Upsert(id string, cab Cab) error
	Delete(id string) error
	Within(center Location, radius float64, limit uint64) ([]Cab, error)
	DeleteAll() error
}

var (
	ErrorNotFound = errors.New("Not found")
	ErrorBadParam = errors.New("Bad parameter")
)
