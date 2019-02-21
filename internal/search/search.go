package search

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

# tag: test
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

// Search options
// look for URI in object space
// look for URI in subject space
// look for term in object literal (but which one?)  (useless?)

// ref https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values

// Search looks for things..  happy?
func Search(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method)
	log.Println(r.URL.Path)
	log.Println(r.Header)

	keys, ok := r.URL.Query()["s"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 's' is missing")
		// return we need params
	}

	key := keys[0] // Query()["key"] will return an array, only want 1
	log.Println("Url Param 's' is: " + string(key))

	repo, err := utils.LDNDBConn()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("test", r)
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Println(q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	bindings := res.Results.Bindings // map[string][]rdf.Term
	check := bindings[0]["s"].Value

	fmt.Println(check)

	//Content-Type: application/ld+json; charset=utf-8; profile="http://www.w3.org/ns/json-ld#expanded"
	w.Header().Add("Content-Type", "application/ld+json; charset=utf-8; profile=\"http://www.w3.org/ns/json-ld#expanded\"")
	jld := "search the data"
	fmt.Fprintf(w, "%s", jld)
}
