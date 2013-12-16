package findcab

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
)

// Runs the http server.  This server offers more control than the standard go's default http server
// in that when a 'true' is sent to the stop channel, the listener is closed to force a clean shutdown.
func RunServer(server *http.Server, stop chan bool) (stopped chan bool) {
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic(err)
	}
	stopped = make(chan bool)

	// This will be set to true if a shutdown signal is received. This allows us to detect
	// if the server stop is intentional or due to some error.
	fromSignal := false

	// The main goroutine where the server listens on the network connection
	go func(fromSignal *bool) {
		// Serve will block until an error (e.g. from shutdown, closed connection) occurs.
		err := server.Serve(listener)
		if !*fromSignal {
			log.Println("Warning: server stops due to error", err)
		}
		stopped <- true
	}(&fromSignal)

	// Another goroutine that listens for signal to close the network connection
	// on shutdown.  This will cause the server.Serve() to return.
	go func(fromSignal *bool) {
		select {
		case <-stop:
			listener.Close()
			*fromSignal = true // Intentially stopped from signal
			return
		}
	}(&fromSignal)
	return
}

// Returns a http server from given service object
// Registration of URL routes to handler functions that will invoke the service's methods to do CRUD.
// Basic marshal/unmarshal of JSON objects for the REST calls also take place here.
func HttpServer(service CabService) *http.Server {

	router := mux.NewRouter()

	// Create / Update Request
	router.Methods("PUT").Path("/cabs/{cabId}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		cabId, err := strconv.ParseUint(params["cabId"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		cab := Cab{}
		err = json.Unmarshal(body, &cab)

		// A quick check to make sure we have id that matches
		if cab.Id == Id(0) {
			cab.Id = Id(cabId) // fill in the missing Id from the URL
		}

		if cab.Id != Id(cabId) {
			http.Error(w, "Cab Id and URL mismatch", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = service.Upsert(cab)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Get Request
	router.Methods("GET").Path("/cabs/{cabId}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cabId, err := strconv.ParseUint(params["cabId"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		cab, err := service.Read(Id(cabId))
		switch err {
		case nil:
			if jsonStr, err2 := json.Marshal(cab); err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			} else {
				w.Write(jsonStr)
				return
			}

		case ErrorNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Query
	router.Methods("GET").Path("/cabs").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println("form", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		longitude, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
		if err != nil {
			log.Println("lng", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		latitude, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
		if err != nil {
			log.Println("lat", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		radius, err := strconv.ParseFloat(r.FormValue("radius"), 64)
		if err != nil {
			log.Println("rad", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		limit := uint64(8)
		if len(r.FormValue("limit")) > 0 {
			limit, _ = strconv.ParseUint(r.FormValue("limit"), 10, 64)
		}

		cabs, err := service.Query(GeoWithin{
			Center: Location{
				Longitude: longitude,
				Latitude:  latitude,
			},
			Radius: radius,
			Unit:   Meters,
			Limit:  int(limit)})
		switch err {
		case nil:
			if jsonStr, err2 := json.Marshal(cabs); err2 != nil {
				http.Error(w, err2.Error(), http.StatusInternalServerError)
				return
			} else {
				w.Write(jsonStr)
				return
			}

		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// Destroy Request
	router.Methods("DELETE").Path("/cabs/{cabId}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cabId, err := strconv.ParseUint(params["cabId"], 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		err = service.Delete(Id(cabId))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Destroy All Request
	router.Methods("DELETE").Path("/cabs").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := service.DeleteAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	return &http.Server{
		Handler: router,
	}
}
