package impl

import (
	"github.com/gyokuro/findcab"
	"log"
)

// Simple implementation of the CabService interface

func SimpleCabService() *simpleCabService {
	return &simpleCabService{
		cabs: make(map[findcab.Id]findcab.Cab),
	}
}

type simpleCabService struct {
	cabs map[findcab.Id]findcab.Cab
}

func (s *simpleCabService) Read(id findcab.Id) (result findcab.Cab, err error) {
	log.Println("id", id)
	var exists bool
	if result, exists = s.cabs[id]; exists {
		return
	}
	err = findcab.ErrorNotFound
	return
}

func (s *simpleCabService) Upsert(id findcab.Id, cab findcab.Cab) (err error) {
	log.Println("Upsert", id, cab)
	s.cabs[id] = cab
	return nil
}

func (s *simpleCabService) Delete(id findcab.Id) (err error) {
	log.Println("Delete", id)
	delete(s.cabs, id)
	return nil
}

func sanitize(q *findcab.GeoWithin) *findcab.GeoWithin {
	if q.Limit == 0 {
		q.Limit = 8
	}
	if q.Unit == 0 {
		q.Unit = findcab.Meters
	}
	return q
}

func (s *simpleCabService) Query(q findcab.GeoWithin) (cabs []findcab.Cab, err error) {
	sanitize(&q)
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
		if len(cabs) == q.Limit {
			return
		}
	}
	return
}

func (s *simpleCabService) DeleteAll() (err error) {
	log.Println("DeleteAll")
	s.cabs = make(map[findcab.Id]findcab.Cab)
	return
}
