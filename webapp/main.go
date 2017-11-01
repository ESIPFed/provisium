package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	dx "lab.esipfed.org/provisium/webapp/dx"
	handlers "lab.esipfed.org/provisium/webapp/handlers"
	kv "lab.esipfed.org/provisium/webapp/kv"
	services "lab.esipfed.org/provisium/webapp/services"
)

// MyServer struct for mux router
type MyServer struct {
	r *mux.Router
}

// Our main is really just route collection....  from here you can see what URLs go
// to what functions.  It adds in a few default header elements and fires up
// the listener.
func main() {
	// DX router (implements our LODish 303 pattern which should be demonstrated here to ensure alignment)
	// TODO: All three patterns go to the same function..  make this one pattern matcher/one line
	dxroute := mux.NewRouter()
	dxroute.HandleFunc("/id/dataset/{ID}", dx.Redirection)            // id -> doc 303 redirection
	dxroute.HandleFunc("/id/dataset/{ID}/provenance", dx.Redirection) // PROV: prov redirection
	dxroute.HandleFunc("/id/dataset/{ID}/pingback", dx.Redirection)   // PROV: pingback for this resource  (would prefer a master /prov or server)
	http.Handle("/id/", dxroute)

	// Data and prov router (LODish)
	dataset := mux.NewRouter()
	dataset.HandleFunc("/doc/dataset/{ID}", handlers.RenderLP)              // PROV: test cast with Void..  would need to generalize
	dataset.HandleFunc("/doc/dataset/{ID}/provenance", handlers.RenderProv) // PROV: test cast with Void..  would need to generalize
	dataset.HandleFunc("/doc/dataset/{ID}/pingback", handlers.ProvPingback) // PROV: pingback for this resource  (would prefer a master /prov or server)
	http.Handle("/doc/", dataset)

	// Catalog router
	catalog := mux.NewRouter()
	catalog.HandleFunc("/catalog/listing", handlers.CatalogListing) // PROV: test cast with Void..  would need to generalize
	http.Handle("/catalog/", catalog)

	// TODO  make all this  :)
	// Section 4.2 https://www.w3.org/TR/2013/NOTE-prov-aq-20130430/#direct-http-query-service-invocation
	// Services router
	sr := mux.NewRouter()
	// services.HandleFunc("/api/v1/prov/service/{CALLSTRING}", services.ProvAQService)
	sr.HandleFunc("/api/v1/event/{ID}", services.ProvEventInfo)
	sr.HandleFunc("/api/v1/docgraph/{ID}", services.DocGraph)
	sr.HandleFunc("/api/v1/docevent/{ID}", services.DocReport)
	http.Handle("/api/", sr)

	// Index router, handle our main page uniquely...   may want to do some things with this eventulay
	root := mux.NewRouter()
	root.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", &MyServer{root})

	// Static router for images, css, js, etc...  (assets)
	static := mux.NewRouter()
	static.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/static/", &MyServer{static})

	// Need a good 404 router

	// Init the KV store to ensure all buckets are ready....
	err := kv.InitKV()
	if err != nil {
		log.Fatal(err) // fatal since if buckets are not ready we can't go play....
	}

	// Start the server...
	log.Printf("About to listen on 9900. Go to http://127.0.0.1:9900/")
	err = http.ListenAndServe(":9900", nil)
	if err != nil {
		log.Fatal(err) // fatal if we can't serve, just go home...
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

	// Let the Gorilla work
	s.r.ServeHTTP(rw, req)
}
