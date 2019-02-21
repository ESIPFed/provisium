package prov

import (
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/knakk/sparql"
	"lab.esipfed.org/provisium/internal/utils"
)

const queries = `
# Comments are ignored, except those tagging a query.

# tag: pingbacks
prefix schema: <http://schema.org/>
prefix bds: <http://www.bigdata.com/rdf/search#>
select *
from <http://provisium.io/graph/id/bhm5hnqu6s70r2q7v71g>
where {
   ?s ?p ?o
}
`

// SPO is a struct to hold boolean check on a resource
type SPO struct {
	S string
	P string
	O string
}

// Doc return html rep of the doc with some examples of how it can be referenced
func Doc(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	log.Println(r.URL.Path)
	log.Println(r.Header)
	log.Println(r.Header["Content-type"])

	//  get your ID..   collect the calls into structus
	//  this is the only call to render an HTML view
	//  Assemble into a single struct (the info) and also return a JSON view
	//  elements:   has_prov   has_query, has_pingback

	//Content-Type: application/ld+json; charset=utf-8; profile="http://www.w3.org/ns/json-ld#expanded"
	w.Header().Add("Content-Type", "test/plain; charset=utf-8; profile=\"http://www.w3.org/ns/json-ld#expanded\"")
	jld := "An HTML rendered view here"
	fmt.Fprintf(w, "%s", jld)
}

func doSearch(uri string) {
	repo, err := utils.LDNDBConn()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("pingbacks", uri)
	if err != nil {
		log.Printf("%s\n", err)
	}

	// log.Println(q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	bindings := res.Results.Bindings // map[string][]rdf.Term
	check := bindings[0]["s"].Value

	fmt.Println(check)

}
