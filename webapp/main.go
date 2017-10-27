package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	dx "lab.esipfed.org/provisium/webapp/dx"
	handlers "lab.esipfed.org/provisium/webapp/handlers"
)

// MyServer struct for mux router
type MyServer struct {
	r *mux.Router
}

// func rootHandler(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	w.WriteHeader(http.StatusOK)
// 	data, err := ioutil.ReadFile("./static/index.html")
// 	if err != nil {
// 		panic(err)
// 	}
// 	w.Header().Set("Content-Length", fmt.Sprint(len(data)))
// 	fmt.Fprint(w, string(data))
// }

func main() {
	// Recall /id is going to be our dx..   all items that come in with that will be looked up and 303'd
	// Example URL:  http://opencoredata.org/id/dataset/c2d80e2a-cc30-430c-b0bd-cee9092688e3
	dxroute := mux.NewRouter()
	dxroute.HandleFunc("/id/dataset/{ID}", dx.Redirection)
	dxroute.HandleFunc("/id/dataset/{ID}/provenance", dx.Redirection) // PROV: prov redirection
	dxroute.HandleFunc("/id/dataset/{ID}/pingback", dx.Redirection)   // PROV: pingback for this resource  (would prefer a master /prov or server)
	http.Handle("/id/", dxroute)

	// Some early Prov Pingback work here...
	dataset := mux.NewRouter()
	dataset.HandleFunc("/doc/dataset/{ID}", handlers.RenderLP)              // PROV: test cast with Void..  would need to generalize
	dataset.HandleFunc("/doc/dataset/{ID}/provenance", handlers.RenderProv) // PROV: test cast with Void..  would need to generalize
	dataset.HandleFunc("/doc/dataset/{ID}/pingback", handlers.ProvPingback) // PROV: pingback for this resource  (would prefer a master /prov or server)
	http.Handle("/doc/", dataset)

	// Catalog listing
	catalog := mux.NewRouter()
	catalog.HandleFunc("/catalog/listing", handlers.CatalogListing) // PROV: test cast with Void..  would need to generalize
	http.Handle("/catalog/", catalog)

	// Index handler
	parking := mux.NewRouter()
	parking.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", &MyServer{parking})

	// Need a good 404 handler

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
