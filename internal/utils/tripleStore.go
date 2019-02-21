package utils

import (
	"log"
	"time"

	"github.com/knakk/sparql"
)

var SPARQLurl string

// LDNDBConn creates a connection to the triple store
func LDNDBConn() (*sparql.Repo, error) {
	repo, err := sparql.NewRepo(SPARQLurl,
		sparql.Timeout(time.Millisecond*15000),
	)
	if err != nil {
		log.Printf("%s\n", err)
	}
	return repo, err
}
