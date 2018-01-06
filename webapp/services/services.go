package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	kv "lab.esipfed.org/provisium/webapp/kv"
)

// PingBackCatalog returns prov for a given document ID as a service defined in the W3C Prov-aq NOTE
// This is RDF prov..  not a URI list..
// http://www.example.com/provenance/service?target={+uri}{&steps}
// as noted in the docs.  Optional steps can be, for us, scope [curated, community]
func PingBackCatalog(w http.ResponseWriter, r *http.Request) {
	ID := r.FormValue("target")
	urlElements := strings.Split(ID, "/")

	events, err := kv.GetProvLog(urlElements[len(urlElements)-1])
	if err != nil {
		log.Println("error getting info..  need to 3** a return and log")
	}

	var buffer bytes.Buffer // a buffere for what we are about to collect
	for k, _ := range events {
		content, _, err := kv.GetProvDetails(k)
		if err != nil {
			log.Println("error getting info..  need to 3** a return and log")
		}
		buffer.WriteString(content + "\n")

	}

	fmt.Println(buffer.String()) // make as small local function to return a unique set of strings

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(buffer.String()))
}

// PingBackEvents will get the information we know about a prov ID
// and return it as something like RDF (JSON-LD, turtle..  not sure)
func PingBackEvents(w http.ResponseWriter, r *http.Request) {
	// I really only deal with the UUID, so I need to strip this off.  It would likely
	// be better to do as the NOTE notes (sigh) and regsiter full qualified addresses
	// This is bad form with "logic in code", and I should pass what I expect to get
	ID := r.FormValue("target")
	urlElements := strings.Split(ID, "/")

	events, err := kv.GetProvLog(urlElements[len(urlElements)-1])
	if err != nil {
		log.Println("error getting info..  need to 3** a return and log")
	}

	jsonString, err := json.MarshalIndent(events, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

// PingBackEvents will get the information we know about a prov ID
// and return it as something like RDF (JSON-LD, turtle..  not sure)
func PingBackContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["ID"]

	content, _, err := kv.GetProvDetails(ID)
	if err != nil {
		log.Println("error getting info..  need to 3** a return and log")
	}

	// I could get the content type..  it is in the KV and returned above, but I ignore it
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(content))
}
