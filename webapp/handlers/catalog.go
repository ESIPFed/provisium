package handlers

import (
	"html/template"
	"log"
	"net/http"

	kv "lab.esipfed.org/provisium/webapp/kv"
)

// CatalogListing displays the catalog of files....
func CatalogListing(w http.ResponseWriter, r *http.Request) {

	// get an array of file IDs

	// mock := [2]string{"id1", "id2"}
	mock := kv.GetDocIDs()

	ht, err := template.New("Landing page template").ParseFiles("templates/catalog.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", mock) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}
