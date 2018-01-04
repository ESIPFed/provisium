package handlers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/knakk/rdf"
	"lab.esipfed.org/provohash/kv"
)

// PostProv takes a POST call to put in the prov with
func PostProv(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 100000) // Max 100KB of prov

	contentType := r.Header.Get("Content-Type")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	// TODO  check the prov to see if it's valid RDF
	validRDF := isValidRDF(string(body))
	if !validRDF {
		log.Print("Not valid RDF..  need to deal with this")
	}

	// Get the prov ID
	md5string, err := kv.NewProvHash(string(body), r.RemoteAddr, contentType)

	log.Print("Log a prov creation event")

	if md5string == "" {
		w.Write([]byte(fmt.Sprintf("I will send an error code too. . but you sent in existing text")))
	} else {
		w.Write([]byte(fmt.Sprintf("http://%s/id/hash/%s", r.Host, md5string)))
	}
}

func isValidRDF(prov string) bool {
	validity := false

	var inoutFormat rdf.Format
	inoutFormat = rdf.Turtle // FormatNT
	dec := rdf.NewTripleDecoder(strings.NewReader(prov), inoutFormat)
	_, err := dec.DecodeAll()
	if err == nil {
		validity = true
	}

	return validity
}

// GetProv returns a landing page for the prov
func GetProv(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["ID"]
	prov, _, err := kv.GetProvDetails(hash)
	if err != nil {
		log.Print("Need to deal with error with correct code")
	}

	// Get template and display landing page meshed with metadata
	ht, err := template.New("Landing page template").ParseFiles("templates/landingPage.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", prov) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
