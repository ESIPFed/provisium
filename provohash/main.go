package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	dx "lab.esipfed.org/provohash/dx"
	handlers "lab.esipfed.org/provohash/handlers"
	kv "lab.esipfed.org/provohash/kv"
	"lab.esipfed.org/provohash/services"
)

// MyServer struct for mux router
type MyServer struct {
	r *mux.Router
}

func main() {
	// DX router (implements our LODish 303 pattern which should be demonstrated here to ensure alignment)
	dxroute := mux.NewRouter()
	dxroute.HandleFunc("/id/hash/{ID}", dx.Redirection) // id -> doc 303 redirection
	http.Handle("/id/", dxroute)

	// Event POST and GET
	dataset := mux.NewRouter()
	dataset.HandleFunc("/doc/newprov", handlers.PostProv)  // PROV: test cast with Void..  would need to generalize
	dataset.HandleFunc("/doc/hash/{ID}", handlers.GetProv) // PROV: test cast with Void..  would need to generalize
	http.Handle("/doc/", dataset)

	// Services
	service := mux.NewRouter()
	service.HandleFunc("/api/v1/prov/{ID}", services.ProvItem) // PROV: test cast with Void..  would need to generalize
	http.Handle("/api/v1/", service)

	// Catalog router  TODO:  Does this need a range and offset?
	catalog := mux.NewRouter()
	catalog.HandleFunc("/catalog/listing", handlers.CatalogListing) // PROV: test cast with Void..  would need to generalize
	http.Handle("/catalog/", catalog)

	// Index router, handle our main page uniquely...   may want to do some things with this eventually
	root := mux.NewRouter()
	root.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", &MyServer{root})

	// Static router for images, css, js, etc...  (assets)
	static := mux.NewRouter()
	static.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/static/", &MyServer{static})

	// Need a good 404 router
	// I keep saying this..  I keep not coding it...   I should do this...

	// Init the KV store to ensure all buckets are ready....
	err := kv.InitKV()
	if err != nil {
		log.Fatal(err) // fatal since if buckets are not ready we can't go play....
	}

	// Start the server...
	log.Printf("About to listen on 9911. Go to http://127.0.0.1:9911/")
	err = http.ListenAndServe(":9911", nil)
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
