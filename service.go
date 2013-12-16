package findcab

import (
	"errors"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type Id uint64

type Cab struct {
	Id        Id      `json:"id" bson:"_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type DistanceUnit int

const (
	Meters DistanceUnit = iota
	Kilometers
	Feet
	Miles
)

type GeoWithin struct {
	Center Location
	Radius float64
	Unit   DistanceUnit
	Limit  int
}

type CabService interface {
	Read(id Id) (Cab, error)
	Upsert(id Id, cab Cab) error
	Delete(id Id) error
	Query(query GeoWithin) ([]Cab, error)
	DeleteAll() error
}

var (
	ErrorNotFound = errors.New("Not found")
	ErrorBadParam = errors.New("Bad parameter")
)
