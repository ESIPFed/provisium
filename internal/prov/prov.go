package prov

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

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

// ID performs a simple content redirection to doc for all cases
// other than content types for RDF, for that it returns the content
func ID(w http.ResponseWriter, r *http.Request) {
	// TODO   if RDF..   return the rdf
	// if anything else redirect to web view
	log.Println(r.Method)
	log.Println(r.URL.Path)
	log.Println(r.Header)
	log.Println(r.Header["Content-Type"])

	// TODO  need to  see if this is even a valid
	// URL to request on.

	p := r.URL.Path
	if utils.Contains(r.Header["Content-Type"], "application/ld+json") {
		w.Header().Add("Content-Type", "application/ld+json; charset=utf-8; profile=\"http://www.w3.org/ns/json-ld#expanded\"")
		jld := TestJSON()
		fmt.Fprintf(w, "%s", jld)
	} else {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8; profile=\"http://www.w3.org/ns/json-ld#expanded\"")
		p = strings.Replace(p, "/id/", "/doc/", -1)
		fmt.Printf("http://" + r.Host + p)
		http.Redirect(w, r, "http://"+r.Host+p, http.StatusMovedPermanently)
	}
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
