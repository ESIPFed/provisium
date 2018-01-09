package services

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"lab.esipfed.org/provisium/provohash/kv"
)

// ProvItem will get the prov by hash
// and return it as something like RDF
func ProvItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["ID"]

	entry, contentType, err := kv.GetProvDetails(ID)
	if err != nil {
		log.Println("error getting info..  need to 3** a return and log")
	}

	w.Header().Set("Content-Type", contentType)
	w.Write([]byte(entry))
}
