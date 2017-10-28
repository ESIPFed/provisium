package services

import "net/http"

// ProvEventInfo will get the information we know about a prov ID
// and return it as something like RDF (JSON-LD, turtle..  not sure)
func ProvEventInfo(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Whoa there..   you found the edge of our expanding universe"))

}
