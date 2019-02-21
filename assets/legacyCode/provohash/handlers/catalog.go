package handlers

import (
	"html/template"
	"log"
	"net/http"

	"lab.esipfed.org/provisium/provohash/kv"
)

// CatalogListing displays the catalog of files....
func CatalogListing(w http.ResponseWriter, r *http.Request) {

	// get an array of file IDs

	// mock := [2]string{"id1", "id2"}
	hashes := kv.GetListing()

	ht, err := template.New("Landing page template").ParseFiles("templates/catalogPage.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", hashes) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
