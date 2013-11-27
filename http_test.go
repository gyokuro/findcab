package findcab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	_ "log"
	"net/http"
	"strconv"
	"testing"
)

var client = &http.Client{}

// A simple implementation that will verify the functionality of the http
// wrapper
type testService struct {
	// record the parameters passed to the api calls
	id     string
	center Location
	radius float64
	limit  uint64
	cab    Cab

	// which method is called?
	calledRead, calledUpsert, calledWithin, calledDelete, calledDeleteAll bool
}

func (ts *testService) Read(id string) (cab Cab, err error) {
	ts.calledRead = true
	ts.id = id
	return
}

func (ts *testService) Upsert(id string, cab Cab) (err error) {
	ts.calledUpsert = true
	ts.id = id
	ts.cab = cab
	return nil
}

func (ts *testService) Delete(id string) (err error) {
	ts.calledDelete = true
	ts.id = id
	return
}

func (ts *testService) Within(center Location, radius float64, limit uint64) (cabs []Cab, err error) {
	ts.calledWithin = true
	ts.center = center
	ts.radius = radius
	ts.limit = limit
	return
}

func (ts *testService) DeleteAll() (err error) {
	ts.calledDeleteAll = true
	return
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func runServer(port int) (service *testService, stop chan bool, stopped chan bool) {
	ts := testService{}
	service = &ts
	httpServer := HttpServer(&ts)
	httpServer.Addr = ":" + strconv.Itoa(port)
	stop = make(chan bool)
	stopped = RunServer(httpServer, stop)
	return
}

func TestHttpCreateUpdate(test *testing.T) {
	service, stop, stopped := runServer(8181)

	cab := Cab{
		Id:        "1234",
		Latitude:  10.,
		Longitude: 100.,
	}
	json, err := json.Marshal(cab)
	check(err)

	req, err := http.NewRequest("PUT", "http://localhost:8181/cabs/1234", bytes.NewBuffer(json))
	check(err)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	check(err)

	stop <- true
	<-stopped

	if resp.StatusCode != 200 {
		test.Error("Expect 200", resp)
	}

	expected := testService{
		calledUpsert: true,
		id:           "1234",
		cab:          cab,
	}

	if *service != expected {
		test.Error("Upsert failed", expected, *service)
	}
}

func TestHttpGet(test *testing.T) {
	service, stop, stopped := runServer(8182)

	req, err := http.NewRequest("GET", "http://localhost:8182/cabs/1234", nil)
	check(err)
	resp, err := client.Do(req)
	check(err)

	stop <- true
	<-stopped

	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	cab := Cab{}
	err = json.Unmarshal(body, &cab)
	check(err)

	expected := testService{
		calledRead: true,
		id:         "1234",
	}

	if *service != expected {
		test.Error("Read failed", expected, *service)
	}
}

func TestHttpQuery(test *testing.T) {
	port := 8183
	service, stop, stopped := runServer(port)

	lat := 5.5
	lng := 15.15
	radius := 1000.0
	limit := uint64(8)

	url := fmt.Sprintf("http://localhost:%d/cabs?latitude=%f&longitude=%f&radius=%f&limit=%d",
		port, lat, lng, radius, limit)
	req, err := http.NewRequest("GET", url, nil)
	check(err)

	resp, err := client.Do(req)
	check(err)

	stop <- true
	<-stopped

	if resp.StatusCode != 200 {
		test.Error("Expect 200", resp)
	}

	expected := testService{
		calledWithin: true,
		center: Location{
			Latitude:  lat,
			Longitude: lng,
		},
		radius: radius,
		limit:  limit,
	}

	if *service != expected {
		test.Error("Query failed", expected, *service)
	}
}

func TestHttpDestroy(test *testing.T) {
	port := 8184
	service, stop, stopped := runServer(port)

	id := "12345"
	url := fmt.Sprintf("http://localhost:%d/cabs/%s", port, id)
	req, err := http.NewRequest("DELETE", url, nil)
	check(err)

	resp, err := client.Do(req)
	check(err)

	stop <- true
	<-stopped

	if resp.StatusCode != 200 {
		test.Error("Expect 200", resp)
	}

	expected := testService{
		calledDelete: true,
		id:           id,
	}

	if *service != expected {
		test.Error("Delete failed", expected, *service)
	}
}

func TestHttpDestroyAll(test *testing.T) {
	port := 8185
	service, stop, stopped := runServer(port)

	url := fmt.Sprintf("http://localhost:%d/cabs", port)
	req, err := http.NewRequest("DELETE", url, nil)
	check(err)

	resp, err := client.Do(req)
	check(err)

	stop <- true
	<-stopped

	if resp.StatusCode != 200 {
		test.Error("Expect 200", resp)
	}

	expected := testService{
		calledDeleteAll: true,
	}

	if *service != expected {
		test.Error("DeleteAll failed", expected, *service)
	}
}
