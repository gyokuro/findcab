package impl

import (
	"github.com/gyokuro/findcab"
	"log"
)

// Dummy implementation of the CabService interface

func DummyCabService() *dummyCabService {
	return &dummyCabService{
		cabs: make(map[string]findcab.Cab),
	}
}

type dummyCabService struct {
	cabs map[string]findcab.Cab
}

func (s *dummyCabService) Read(id string) (result findcab.Cab, err error) {
	log.Println("id", id)
	var exists bool
	if result, exists = s.cabs[id]; exists {
		return
	}
	err = findcab.ErrorNotFound
	return
}

func (s *dummyCabService) Upsert(id string, cab findcab.Cab) (err error) {
	log.Println("Upsert", id, cab)
	s.cabs[id] = cab
	return nil
}

func (s *dummyCabService) Delete(id string) (err error) {
	log.Println("Delete", id)
	delete(s.cabs, id)
	return nil
}

func (s *dummyCabService) Query(q findcab.GeoWithin) (cabs []findcab.Cab, err error) {
	log.Println("Query", q)
	cabs = make([]findcab.Cab, 0)
	for _, cab := range s.cabs {
		distance := Haversine(q.Center, findcab.Location{
			Latitude:  cab.Latitude,
			Longitude: cab.Longitude,
		}, q.Unit)
		if distance <= q.Radius {
			cabs = append(cabs, cab)
		}
	}
	return
}

func (s *dummyCabService) DeleteAll() (err error) {
	log.Println("DeleteAll")
	s.cabs = make(map[string]findcab.Cab)
	return
}
