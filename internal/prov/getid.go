package prov

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

const iq = `
# Comments are ignored, except those tagging a query.
# tag: byid
CONSTRUCT { ?s ?p ?o } WHERE { GRAPH <http://provisium.io/prov/id/{{.}}> { ?s ?p ?o } . }
`

// ID performs a simple content redirection to doc for all cases
// other than content types for RDF, for that it returns the content
func ID(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	pa := strings.Split(p, "/")
	xid := pa[len(pa)-1]
	if utils.Contains(r.Header["Content-Type"], "application/n-triples") {
		w.Header().Add("Content-Type", "application/n-triples; charset=utf-8")
		jld := idSearch(xid)
		fmt.Fprintf(w, "%s", jld)
	} else {
		w.Header().Add("Content-Type", "text/plain; charset=utf-8; profile=\"http://www.w3.org/ns/json-ld#expanded\"")
		p = strings.Replace(p, "/id/", "/doc/", -1)
		fmt.Printf("http://" + r.Host + p)
		http.Redirect(w, r, "http://"+r.Host+p, http.StatusMovedPermanently)
	}
}

func idSearch(xid string) string {
	repo, err := utils.LDNDBConn()
	if err != nil {
		log.Printf("%s\n", err)
	}

	f := bytes.NewBufferString(iq)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("byid", xid)
	if err != nil {
		log.Printf("%s\n", err)
	}

	//log.Println(q)

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
