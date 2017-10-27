package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/knakk/rdf"
	kv "lab.esipfed.org/provisium/webapp/kv"
)

// RenderLP displays the RDF resource and adds a prov pingback entry
func RenderLP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["ID"]

	// get metadata for this document
	mock, err := kv.GetResMetaData(ID)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(r.URL.Path[5:])
	fmt.Println(r.URL.Path[1:])
	// Note the HACK in the next line..  this is just an ALPHA..  (even so this sucks)
	// TODO..   think about issues of 303 here with /id/ and /doc/ since that could becomine an issue...
	// TODO..   it's hard to expect community clients to read and address 303?
	linkProv := fmt.Sprintf("<http://%s/doc/%s/provenance>; rel=\"http://www.w3.org/ns/prov#has_provenance\"", r.Host, r.URL.Path[5:]) // use r.Host so we don't hardcode in
	linkPB := fmt.Sprintf("<http://%s/doc/%s/pingback>; rel=\"http://www.w3.org/ns/prov#pingbck\"", r.Host, r.URL.Path[5:])
	w.Header().Add("Link", linkProv)
	w.Header().Add("Link", linkPB)
	// w.Header().Set("Content-type", "text/plain")

	// Get template and display landing page meshed with metadata

	// http.ServeFile(w, r, fmt.Sprintf("./static/%s", r.URL.Path[1:]))

	ht, err := template.New("Landing page template").ParseFiles("templates/landingPage.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", mock) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
	if err != nil {
		log.Printf("htemplate execution failed: %s", err)
	}

}

// RenderProv shows the prov of a resource, it's just a dummy function now.....
// TODO  need to have this actually get some PROV  :)
func RenderProv(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(getProvRecord()))
}

// ProvPingback Handles the PROV pingback on a resource
func ProvPingback(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
	}

	fmt.Printf("Prov for %s\n", r.URL.Path[1:])
	fmt.Println(string(body))
	// TODO
	// 1) validate this this  (400 if not)
	// 2) store this this to KV store
	// 3) Rolling to the master triple store..
	w.WriteHeader(http.StatusNoContent)
}

// getProvRecord an un-exported function to generate MOCK prov data for testing
// This is a test function only
func getProvRecord() string {

	tr := []rdf.Triple{}

	// Add in
	newsub, _ := rdf.NewIRI("http://foo.org/thisSample")
	newpred1, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	newobj1, _ := rdf.NewIRI("http://www.w3.org/ns/prov#entity")
	tr = append(tr, rdf.Triple{Subj: newsub, Pred: newpred1, Obj: newobj1})

	ga, _ := rdf.NewIRI("http://opencoredata.org/org") // ?
	newpred2, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	newobj2, _ := rdf.NewIRI("http://www.w3.org/ns/prov#org")
	tr = append(tr, rdf.Triple{Subj: ga, Pred: newpred2, Obj: newobj2})

	bn, _ := rdf.NewBlank("bn1")

	newpred3, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	newobj3, _ := rdf.NewIRI("http://www.w3.org/ns/prov#attribution")
	tr = append(tr, rdf.Triple{Subj: bn, Pred: newpred3, Obj: newobj3})

	newpred4, _ := rdf.NewIRI("http://www.w3.org/ns/prov#agent")
	tr = append(tr, rdf.Triple{Subj: bn, Pred: newpred4, Obj: ga})

	newpred5, _ := rdf.NewIRI("http://www.w3.org/ns/prov#hadRole")
	newobj5, _ := rdf.NewIRI("http://www.aurole.org/Publisher")
	tr = append(tr, rdf.Triple{Subj: bn, Pred: newpred5, Obj: newobj5})

	newpred6, _ := rdf.NewIRI("http://www.w3.org/ns/prov#qualifiedAttribution")
	tr = append(tr, rdf.Triple{Subj: newsub, Pred: newpred6, Obj: bn})

	newpred7, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type")
	newobj7, _ := rdf.NewIRI("http://www.w3.org/ns/prov#Agent")
	tr = append(tr, rdf.Triple{Subj: ga, Pred: newpred7, Obj: newobj7})

	newpred8, _ := rdf.NewIRI("http://www.w3.org/ns/prov#wasAttributedTo")
	tr = append(tr, rdf.Triple{Subj: newsub, Pred: newpred8, Obj: bn})

	newpred9, _ := rdf.NewIRI("http://www.w3.org/1999/02/22-rdf-syntax-ns#label")
	newobj9, _ := rdf.NewLiteral("Geoscience Australia")
	tr = append(tr, rdf.Triple{Subj: ga, Pred: newpred9, Obj: newobj9})

	fmt.Println(tr)

	var inoutFormat rdf.Format
	inoutFormat = rdf.Turtle //NTriples

	// Create a buffer io writer
	var b bytes.Buffer
	foo := bufio.NewWriter(&b)

	enc := rdf.NewTripleEncoder(foo, inoutFormat)
	err := enc.EncodeAll(tr)
	err = enc.Close()
	if err != nil {
		log.Println(err)
	}

	foo.Flush()
	return string(b.Bytes())

}
