package search

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/knakk/rdf"
	"github.com/knakk/sparql"
	"lab.esipfed.org/provisium/internal/utils"
)

const sq = `
# Comments are ignored, except those tagging a query.
# tag: byid
CONSTRUCT { ?s ?p ?o } WHERE { GRAPH <http://provisium.io/prov/id/{{.}}> { ?s ?p ?o } . }
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
	keys, ok := r.URL.Query()["s"]
	if !ok || len(keys[0]) < 1 {
		log.Println("Url Param 's' is missing")
	}

	key := keys[0] // Query()["key"] will return an array, only want 1
	log.Println("Url Param 's' is: " + string(key))

	//Content-Type: application/ld+json; charset=utf-8; profile="http://www.w3.org/ns/json-ld#expanded"
	w.Header().Add("Content-Type", "application/n-triples; charset=utf-8")
	jld := idSearch(string(key))
	fmt.Fprintf(w, "%s", jld)
}

func idSearch(xid string) string {
	repo, err := utils.LDNDBConn()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(sq)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("byid", xid)
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Println(q)

	res, err := repo.Construct(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	var b strings.Builder
	for i := range res {
		fmt.Fprintf(&b, "%s", res[i].Serialize(rdf.NTriples))
	}

	return b.String()
}
