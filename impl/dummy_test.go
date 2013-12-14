package impl

import (
	"github.com/gyokuro/findcab"
	"testing"
)

var (
	service = DummyCabService()

	cab1 = findcab.Cab{
		Id:        "1",
		Longitude: -77.037852,
		Latitude:  38.898556,
	}

	cab2 = findcab.Cab{
		Id:        "2",
		Longitude: -77.037852,
		Latitude:  48.898557,
	}
)

func testUpsert(service findcab.CabService, test *testing.T) {
	err := service.Upsert(cab1.Id, cab1)
	if err != nil {
		test.Error("Got error", err)
	}

	err = service.Upsert(cab2.Id, cab2)
	if err != nil {
		test.Error("Got error", err)
	}
}

func testQueryMatch(service findcab.CabService, test *testing.T) {
	q := findcab.GeoWithin{
		Center: findcab.Location{
			Longitude: -77.043934,
			Latitude:  38.897147,
		},
		Radius: 1000.,
		Unit:   findcab.Meters,
	}

	cabs, err := service.Query(q)
	if err != nil {
		test.Error("Got error", err)
	}

	if cabs == nil || len(cabs) != 1 {
		test.Error("Expect one match", cabs)
	}

	if cabs[0] != cab1 {
		test.Error("Expecting", cab1, "got", cabs[0])
	}
}

func testQueryNoMatch(service findcab.CabService, test *testing.T) {
	q := findcab.GeoWithin{
		Center: findcab.Location{
			Longitude: -77.043934,
			Latitude:  38.897147,
		},
		Radius: 100.,
		Unit:   findcab.Meters,
	}

	cabs, err := service.Query(q)
	if err != nil {
		test.Error("Got error", err)
	}

	if cabs == nil || len(cabs) != 0 {
		test.Error("Expect zero match with empty array", cabs)
	}
}

func testDelete(service findcab.CabService, test *testing.T) {
	err := service.Delete(cab2.Id)
	if err != nil {
		test.Error("Got error", err)
	}
}

func testDeleteAll(service findcab.CabService, test *testing.T) {
	err := service.DeleteAll()
	if err != nil {
		test.Error("Got error", err)
	}
}

func TestDummyUpsert(test *testing.T) {
	testUpsert(service, test)
}

func TestDummyQuery(test *testing.T) {
	testQueryMatch(service, test)
	testQueryNoMatch(service, test)
}

func TestDummyDelete(test *testing.T) {
	testDelete(service, test)
	testQueryMatch(service, test)
	testQueryNoMatch(service, test)
}

func TestDummyDeleteAll(test *testing.T) {
	testDeleteAll(service, test)
	testQueryNoMatch(service, test)
}
