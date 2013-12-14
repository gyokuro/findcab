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

// Returns a channel that can be blocked on to get notification of the server's having stopped.
func RunServer(server *http.Server, stop chan bool) (stopped chan bool) {
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		panic(err)
	}
	stopped = make(chan bool)

	// The main goroutine where the server listens on the network connection
	go func() {
		err := server.Serve(listener)
		log.Println("Stopped http server", err)
		stopped <- true
	}()

	// Another goroutine that listens for signal to close the network connection
	// on shutdown.  This will cause the server.Serve() to return.
	go func() {
		select {
		case <-stop:
			listener.Close()
			return
		}
	}()
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
		cabId := params["cabId"]
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		cab := Cab{}
		err = json.Unmarshal(body, &cab)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		err = service.Upsert(cabId, cab)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Get Request
	router.Methods("GET").Path("/cabs/{cabId}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cabId := params["cabId"]
		cab, err := service.Read(cabId)
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
			Limit:  limit})
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
		cabId := params["cabId"]
		err := service.Delete(cabId)
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
