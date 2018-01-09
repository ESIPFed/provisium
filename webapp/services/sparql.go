package services

import (
	"bytes"
	"log"
	"time"

	sparql "github.com/knakk/sparql"
)

const queries = `
# Comments are ignored, except those tagging a query.

# tag: provrecord 
SELECT *
WHERE 
{ 
  <{{.URI}}>  ?p ?o .
}

# tag: provDS
select ?type ?qa  ?at ?qatype ?attype  ?agent ?role 
where {
	<{{.URI}}> a ?type  .
	<{{.URI}}>  <http://www.w3.org/ns/prov#qualifiedAttribution> ?qa  .
	<{{.URI}}>  <http://www.w3.org/ns/prov#wasAttributedTo> ?at .
  ?qa a ?qatype .
  ?at a ?attype .
  ?qa  <http://www.w3.org/ns/prov#agent>  ?agent .
  ?at  <http://www.w3.org/ns/prov#hadRole> ?role .
}
`

// GetProv pulls prov for a resource
func GetProv(uri string) *sparql.Results {
	repo, err := getSPARQL()

	f := bytes.NewBufferString(queries)
	bank := sparql.LoadBank(f)

	q, err := bank.Prepare("provDS", struct{ URI string }{uri})
	if err != nil {
		log.Printf("%s\n", err)
	}

	log.Println(q)

	res, err := repo.Query(q)
	if err != nil {
		log.Printf("%s\n", err)
	}

	return res
}

func getSPARQL() (*sparql.Repo, error) {
	repo, err := sparql.NewRepo("http://localhost:9999/blazegraph/namespace/prov/sparql",
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return repo, err
}
