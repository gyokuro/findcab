package findcab

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

	_, err = client.Do(req)
	check(err)

	stop <- true
	<-stopped

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

	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	cab := Cab{}
	err = json.Unmarshal(body, &cab)
	check(err)

	test.Log("Response:", body, cab)

	stop <- true
	<-stopped

	expected := testService{
		calledRead: true,
		id:         "1234",
	}

	if *service != expected {
		test.Error("Read failed", expected, *service)
	}
}

func TestHttpQuery(test *testing.T) {

}

func TestHttpDestroy(test *testing.T) {

}

func TestHttpDestroyAll(test *testing.T) {

}
