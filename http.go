package findcab

import (
	_ "encoding/json"
	"github.com/gorilla/mux"
	_ "io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Location struct {
	Latitude  float64
	Longitude float64
}

type Cab struct {
	Id        string
	Latitude  float64
	Longitude float64
}

type CabService interface {
	Read(id string) (Cab, error)
	Upsert(cab Cab) error
	Delete(id string) error
	Within(center Location, radius float64, limit uint64) ([]Cab, error)
	ReadAll() ([]Cab, error)
	DeleteAll() error
}

func HttpServer(service CabService) *http.Server {

	router := mux.NewRouter()
	router.Methods("GET").Path("/cabs/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cabs, err := service.ReadAll()
		log.Println("found cabs", cabs, err)
	})
	router.Methods("GET").Path("/cabs").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		longitude, err := strconv.ParseFloat(r.FormValue("longitude"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		latitude, err := strconv.ParseFloat(r.FormValue("latitude"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		radius, err := strconv.ParseFloat(r.FormValue("radius"), 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		limit := uint64(8)
		if len(r.FormValue("limit")) > 0 {
			limit, _ = strconv.ParseUint(r.FormValue("limit"), 10, 64)
		}

		cabs, err := service.Within(Location{
			Longitude: longitude,
			Latitude:  latitude,
		}, radius, limit)
		log.Println("found cabs", cabs, err)
	})

	router.Methods("GET").Path("/cabs/{cabId}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cabId := params["cabId"]
		log.Println("GET", cabId)
		cab, err := service.Read(cabId)
		log.Println("Found", cab, err)
	})

	router.Methods("PUT").Path("/cabs/{cabId}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cabId := params["cabId"]
		log.Println("PUT", cabId)
		// TODO - parse bod
	})

	router.Methods("DELETE").Path("/cabs/{cabId}").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		cabId := params["cabId"]
		log.Println("DELETE", cabId)
		// TODO - parse bod
	})

	router.Methods("DELETE").Path("/cabs").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("DELETE ALL")
		err := service.DeleteAll()
		log.Println("DELETED", err)
	})

	return &http.Server{
		Handler: router,
	}
}
