package services

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	kv "lab.esipfed.org/provisium/webapp/kv"
)

// ProvEventInfo will get the information we know about a prov ID
// and return it as something like RDF (JSON-LD, turtle..  not sure)
func ProvEventInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["ID"]

	entry, contentType, err := kv.GetProvDetails(ID)
	if err != nil {
		log.Println("error getting info..  need to 3** a return and log")
	}

	w.Header().Set("Content-Type", contentType)
	w.Write([]byte(entry))
}
