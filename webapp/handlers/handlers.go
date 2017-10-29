package handlers

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/knakk/rdf"
	kv "lab.esipfed.org/provisium/webapp/kv"
)

type PageData struct {
	SchemaOrg string
	EventLog  map[string]string
	ProvRDF   string
}

// RenderLP displays the RDF resource and adds a prov pingback entry
func RenderLP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["ID"]

	// Get schemaOrg and event log for this document
	so, err := kv.GetResMetaData(ID)
	if err != nil {
		log.Println(err)
	}
	events := kv.GetProvLog(ID)                     // TODO, this should return an error too
	pd := PageData{SchemaOrg: so, EventLog: events} // struct to pass to the page

	// Note the HACK in the next line..  this is just an ALPHA..  (even so this sucks)
	// TODO..   think about issues of 303 here with /id/ and /doc/ since that could becomine an issue...
	// TODO..   it's hard to expect community clients to read and address 303?
	linkProv := fmt.Sprintf("<http://%s/doc/%s/provenance>; rel=\"http://www.w3.org/ns/prov#has_provenance\"", r.Host, r.URL.Path[5:]) // use r.Host so we don't hardcode in
	linkPB := fmt.Sprintf("<http://%s/doc/%s/pingback>; rel=\"http://www.w3.org/ns/prov#pingbck\"", r.Host, r.URL.Path[5:])
	w.Header().Add("Link", linkProv)
	w.Header().Add("Link", linkPB)
	// w.Header().Set("Content-type", "text/plain")

	// Get template and display landing page meshed with metadata
	ht, err := template.New("Landing page template").ParseFiles("templates/landingPage.html") //open and parse a template text file
	if err != nil {
		log.Printf("template parse failed: %s", err)
	}

	err = ht.ExecuteTemplate(w, "T", pd) //substitute fields in the template 't', with values from 'user' and write it out to 'w' which implements io.Writer
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

	// TODO  Body should be a URI list... check for content-type: text/uri-list

	fmt.Printf("Prov for %s\n", r.URL.Path[1:])
	pathElements := strings.Split(r.URL.Path[1:], "/")
	docID := pathElements[2]
	fmt.Println(docID)

	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	var URLError error
	for scanner.Scan() {
		_, err := url.ParseRequestURI(scanner.Text()) // validate this is a URL
		if err != nil {
			URLError = err
		}
		fmt.Printf("URL: %s is valid: %v\n", scanner.Text(), err)
	}

	// TODO..  require content type to be set?  must be URI list or RDF of some form
	// if not error..  4** conent not supported
	contentType := r.Header.Get("Content-Type")
	err = kv.NewProvEvent(docID, string(body), r.RemoteAddr, contentType)
	if err != nil {
		fmt.Println("Error trying to record the uploaded prov")
		URLError = err
	}

	if URLError != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

}

// getProvRecord an un-exported function to generate MOCK prov data for testing
// This is a TEST function only
// GET RID OF THIS ASAP
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
