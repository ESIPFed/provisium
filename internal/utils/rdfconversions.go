package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/knakk/rdf"
	"github.com/piprate/json-gold/ld"
)

// NQToJSONLD converts Quads (or triples) back to JSON-LD  (honors graphs)
func NQToJSONLD(triples string) ([]byte, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	// add the processing mode explicitly if you need JSON-LD 1.1 features
	options.ProcessingMode = ld.JsonLd_1_1

	doc, err := proc.FromRDF(triples, options)
	if err != nil {
		panic(err)
	}

	// ld.PrintDocument("JSON-LD output", doc)
	b, err := json.MarshalIndent(doc, "", " ")

	return b, err
}

// JSONLDToNQ get nquads back from JSON-LD (honor any @graph)
// but return ntriples for JSON-LD with no @graph
func JSONLDToNQ(jsonld string) (string, error) {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")
	// add the processing mode explicitly if you need JSON-LD 1.1 features
	options.ProcessingMode = ld.JsonLd_1_1
	options.Format = "application/n-quads"

	var myInterface interface{}
	err := json.Unmarshal([]byte(jsonld), &myInterface)
	if err != nil {
		log.Println("Error when transforming JSON-LD document to interface:", err)
		return "", err
	}

	triples, err := proc.ToRDF(myInterface, options) // returns triples but toss them, just validating
	if err != nil {
		log.Println("Error when transforming JSON-LD document to RDF:", err)
		return "", err
	}

	return fmt.Sprintf("%v", triples), err
}

// JSONLDToNT to get only ntriples (no graph)
func JSONLDToNT(jsonld string) (string, error) {
	nq, err := JSONLDToNQ(jsonld)
	if err != nil {
		log.Println(err)
		return "", err
	}

	var b bytes.Buffer
	writer := bufio.NewWriter(&b)

	dec := rdf.NewTripleDecoder(strings.NewReader(nq), rdf.NTriples)
	tr, err := dec.DecodeAll()
	if err != nil {
		log.Println(err)
	}

	enc := rdf.NewTripleEncoder(writer, rdf.NTriples)
	err = enc.EncodeAll(tr)
	if err != nil {
		log.Println(err)
	}

	// dec := rdf.NewQuadDecoder(strings.NewReader(nq), rdf.NTriples)
	// tr, err := dec.DecodeAll()
	// if err != nil {
	// 	log.Println(err)
	// }

	// enc := rdf.NewQuadEncoder(writer, rdf.NQuads)
	// err = enc.EncodeAll(tr)
	// if err != nil {
	// 	log.Println(err)
	// }

	writer.Flush()
	return b.String(), err
}
