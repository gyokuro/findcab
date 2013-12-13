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

type DistanceUnit float64

const (
	Meters DistanceUnit = iota
	Feet
)

type WithinQuery struct {
	Center Location
	Radius float64
	Unit   DistanceUnit
	Limit  uint64
}

type CabService interface {
	Read(id string) (Cab, error)
	Upsert(id string, cab Cab) error
	Delete(id string) error
	Query(query WithinQuery) ([]Cab, error)
	DeleteAll() error
}

var (
	ErrorNotFound = errors.New("Not found")
	ErrorBadParam = errors.New("Bad parameter")
)
