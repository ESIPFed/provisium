package prov

import (
	"log"
	"net/http"
	"strings"

	"github.com/alecthomas/template"
)

// Doc return html rep of the doc with some examples of how it can be referenced
func Doc(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	pa := strings.Split(p, "/")
	xid := pa[len(pa)-1]
	//Content-Type: application/ld+json; charset=utf-8; profile="http://www.w3.org/ns/json-ld#expanded"
	w.Header().Add("Content-Type", "test/plain; charset=utf-8; profile=\"http://www.w3.org/ns/json-ld#expanded\"")
	d := idSearch(xid)

	// TODO  make a simple HTML template page and load this and a few inclusion examples up
	// for the users  (ie, how to reference your prov)

	ht, err := template.New("some template").ParseFiles("web/templates/doc.html")
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	// TODO make a struct to pass d and xid back to the page to better build out the example

	err = ht.ExecuteTemplate(w, "T", d)
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

	// fmt.Fprintf(w, "%s", jld)
}
