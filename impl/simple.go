package impl

import (
	"github.com/gyokuro/findcab"
)

// Simple implementation of the CabService interface
// This implementation uses a hashmap and does a O(N) scan of all entries
// when computing the nearest neighbor.
type simpleCabService struct {
	cabs map[findcab.Id]findcab.Cab
}

// Constructor method.  Returns an instance of the simple service
func NewSimpleCabService() *simpleCabService {
	return &simpleCabService{
		cabs: make(map[findcab.Id]findcab.Cab),
	}
}

// Implements CabService
func (s *simpleCabService) Read(id findcab.Id) (result findcab.Cab, err error) {
	var exists bool
	if result, exists = s.cabs[id]; exists {
		return
	}
	err = findcab.ErrorNotFound
	return
}

// Implements CabService
func (s *simpleCabService) Upsert(cab findcab.Cab) (err error) {
	s.cabs[cab.Id] = cab
	return nil
}

// Implements CabService
func (s *simpleCabService) Delete(id findcab.Id) (err error) {
	delete(s.cabs, id)
	return nil
}

// Implements CabService
func (s *simpleCabService) Query(q findcab.GeoWithin) (cabs []findcab.Cab, err error) {
	findcab.Sanitize(&q)
	cabs = make([]findcab.Cab, 0)
	for _, cab := range s.cabs {
		distance := Haversine(q.Center, findcab.Location{
			Latitude:  cab.Latitude,
			Longitude: cab.Longitude,
		}, q.Unit)
		if distance <= q.Radius {
			cabs = append(cabs, cab)
		}
		if len(cabs) == q.Limit {
			return
		}
	}
	return
}

// Implements CabService
func (s *simpleCabService) DeleteAll() (err error) {
	s.cabs = make(map[findcab.Id]findcab.Cab)
	return
}

// Implements CabService
func (s *simpleCabService) Close() {
	// no op
}
