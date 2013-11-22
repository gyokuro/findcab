package findcab

import (
	_ "encoding/json"
	"github.com/gorilla/mux"
	_ "io/ioutil"
	_ "log"
	"net/http"
	"strconv"
)

type Server struct {
	Port int
}

func (f *Server) Http() *http.Server {

	rest := mux.NewRouter()
	rest.HandleFunc("/cab/{cabId}",
		func(w http.ResponseWriter, r *http.Request) {
			params := mux.Vars(r)
			cabId := params["cabId"]
			w.Write([]byte(cabId))
		})

	return &http.Server{
		Addr:    ":" + strconv.Itoa(f.Port),
		Handler: rest,
	}
}
