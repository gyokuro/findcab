package impl

import (
	"github.com/gyokuro/findcab"
	"math"
)

const (
	to_radians       float64 = math.Pi / 180.
	EarthRadiusKm            = 6373.
	EarthRadiusMiles         = 3961.
)

/// converts to radians
func rad(deg float64) float64 {
	return deg * to_radians
}

func sin(deg float64) float64 {
	return math.Sin(rad(deg))
}

func cos(deg float64) float64 {
	return math.Cos(rad(deg))
}

func asin(deg float64) float64 {
	return math.Asin(rad(deg))
}

func atan2(deg1, deg2 float64) float64 {
	return math.Atan2(rad(deg1), rad(deg2))
}

// Compute the haversine distance between two locations expressed in lat/lng.
// Source: http://andrew.hedges.name/experiments/haversine/
func Haversine(l1, l2 findcab.Location, unit findcab.DistanceUnit) float64 {
	dlon := l2.Longitude - l1.Longitude
	dlat := l2.Latitude - l1.Latitude

	a := math.Pow(sin(dlat/2), 2) + cos(l1.Latitude)*cos(l2.Latitude)*math.Pow(sin(dlon/2), 2)
	c := 2 * atan2(math.Sqrt(a), math.Sqrt(1-a))
	switch unit {
	case findcab.Kilometers:
		return c * EarthRadiusKm
	case findcab.Meters:
		return c * EarthRadiusKm * 1000.
	case findcab.Miles:
		return c * EarthRadiusMiles
	case findcab.Feet:
		return c * EarthRadiusMiles * 5280.
	}
	return 0.
}
