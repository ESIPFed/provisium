package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	handlers "lab.esipfed.org/provisium/webapp/handlers"
	"opencoredata.org/ocdWeb/dx"
)

// MyServer struct for mux router
type MyServer struct {
	r *mux.Router
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	data, err := ioutil.ReadFile("./static/index.html")
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
	fmt.Fprint(w, string(data))
}

func main() {
	// root
	parking := mux.NewRouter()
	parking.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", &MyServer{parking})

	// Recall /id is going to be our dx..   all items that come in with that will be looked up and 303'd
	// Example URL:  http://opencoredata.org/id/dataset/c2d80e2a-cc30-430c-b0bd-cee9092688e3
	dxroute := mux.NewRouter()
	dxroute.HandleFunc("/id/graph/{id}", dx.RDFRedirection)
	dxroute.HandleFunc("/id/graph/{id}/provenance", dx.ProvRedirection)   // PROV: prov redirection
	dxroute.HandleFunc("/id/graph/{id}/pingback", dx.PingbackRedirection) // PROV: pingback for this resource  (would prefer a master /prov or server)
	dxroute.HandleFunc("/id/dataset/{UUID}", dx.Redirection)
	dxroute.HandleFunc(`/id/resource/{resourcepath:[a-zA-Z0-9=\-\/]+}`, dx.Redirection)
	http.Handle("/id/", dxroute)

	// Some early Prov Pingback work here...   Deal with void...  (show void..  allow .rdf file downloads)
	dataset := mux.NewRouter()
	dataset.HandleFunc("/doc/dataset/{id}", handlers.RenderWithProvHeader)      // PROV: test cast with Void..  would need to generalize
	dataset.HandleFunc("/doc/dataset/{id}/provenance", handlers.RenderWithProv) // PROV: test cast with Void..  would need to generalize
	dataset.HandleFunc("/doc/dataset/{id}/pingback", handlers.ProvPingback)     // PROV: pingback for this resource  (would prefer a master /prov or server)
	http.Handle("/doc/", dataset)

	// Start the server...
	log.Printf("About to listen on 9900. Go to http://127.0.0.1:9900/")
	err := http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// ref http://stackoverflow.com/questions/12830095/setting-http-headers-in-golang
func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	rw.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	// Stop here if its Preflighted OPTIONS request
	// if req.Method == "OPTIONS" {
	// 	return
	// }

	// Lets Gorilla work
	s.r.ServeHTTP(rw, req)
}

func addDefaultHeaders(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fn(w, r)
	}
}
