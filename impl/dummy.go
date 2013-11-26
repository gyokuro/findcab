package impl

import (
	"github.com/gyokuro/findcab"
	"log"
)

// Dummy implementation of the CabService interface

type DummyCabService struct {
}

var dummyCab = findcab.Cab{
	Id:        "0",
	Latitude:  1.0,
	Longitude: 1.0,
}

func (cs *DummyCabService) Read(id string) (cab findcab.Cab, err error) {
	log.Println("id", id)
	if id == dummyCab.Id {
		cab = dummyCab
	} else {
		err = findcab.ErrorNotFound
	}
	return
}

func (cs *DummyCabService) Upsert(id string, cab findcab.Cab) (err error) {
	log.Println("Upsert", cab)
	return nil
}

func (cs *DummyCabService) Delete(id string) (err error) {
	log.Println("Delete", id)
	return nil
}

func (cs *DummyCabService) Within(center findcab.Location,
	radius float64, limit uint64) (cabs []findcab.Cab, err error) {
	log.Println("Within", center, radius, limit)
	return
}

func (cs *DummyCabService) DeleteAll() (err error) {
	log.Println("DeleteAll")
	return
}
