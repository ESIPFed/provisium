package prov

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/rs/xid"
	"lab.esipfed.org/provisium/internal/utils"
)

// HostProvURI accepts RDF prov to store and host for parties who do not wish
// to manage their own local prov
func HostProvURI(w http.ResponseWriter, r *http.Request) {
	// get the POST body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on body parameter read %v \n", err)
		log.Println(err)
		w.WriteHeader(422)
		fmt.Fprintf(w, "")
		return
	}

	// Check that it is JSON-LD and attempt to convert it to Turtle
	// nq, err := utils.JSONLDToNQ(string(body))
	// if err != nil {
	// 	log.Println(err)
	// }

	nt, err := utils.JSONLDToNT(string(body))
	if err != nil {
		log.Println(err)
	}

	id := xid.New() // get a new UID to for the graph URI
	insert := fmt.Sprintf("INSERT DATA {GRAPH <http://provisium.io/prov/id/%s> { %s }}", id, nt)

	// Ref https://www.w3.org/TR/sparql11-http-rdf-update
	//  INSERT DATA { GRAPH <graph_uri> { .. RDF payload .. } }

	// Try to insert into Jena
	_, err = UpdateCall([]byte(insert))
	if err != nil {
		log.Printf("Error on updatecall %v \n", err)
		log.Println(err)
		w.WriteHeader(422)
		fmt.Fprintf(w, "")
		return
	}

	w.Header().Add("Location", fmt.Sprintf("http://provisium.io/prov/id/%s", id))
	w.WriteHeader(201)
	fmt.Fprintf(w, "")
}

// TestJSON is just a place holder function for some JSON[-ld]
func TestJSON() string {
	j := `{
	"@context": {
	  "ical": "http://www.w3.org/2002/12/cal/ical#",
	  "xsd": "http://www.w3.org/2001/XMLSchema#",
	  "ical:dtstart": {
		"@type": "xsd:dateTime"
	  }
	},
	"@id" : "http://example.com/id/1",
	"ical:summary": "Lady Gaga Concert",
	"ical:location": "New Orleans Arena, New Orleans, Louisiana, USA",
	"ical:dtstart": "2011-04-09T20:00:00Z"
  }`
	return j
}
