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

// DocGraph returns prov for a given document ID
// This is RDF prov..  not a URI list..   so..  we need policy on our URI-list holdings
// also, use the pattern:
// http://www.example.com/provenance/service?target={+uri}{&steps}
// as noted in the docs.  Optional steps can be, for us, scope [curated, community]
func DocGraph(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Whoa there....  you have reached the edge of our ever expanding universe"))
}

// DocReport should return all we know about events on the document.  POSTings,
// local graph, URI lists, community graph obtained, etc...
func DocReport(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Whoa there....  you have reached the edge of our ever expanding universe"))
}
