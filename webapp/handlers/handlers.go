package handers

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/knakk/rdf"
)

// RenderWithProvHeader displays the RDF resource and adds a prov pingback entry
func RenderWithProvHeader(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path[5:])                                                                                                       // Note the HACK in the next line..  this is just an ALPHA..  (even this this sucks)
	linkProv := fmt.Sprintf("<http://%s/id/%s/provenance>; rel=\"http://www.w3.org/ns/prov#has_provenance\"", r.Host, r.URL.Path[5:]) // use r.Host so we don't hardcode in
	linkPB := fmt.Sprintf("<http://%s/rdf/%s/pingback>; rel=\"http://www.w3.org/ns/prov#pingbck\"", r.Host, r.URL.Path[5:])
	w.Header().Add("Link", linkProv)
	w.Header().Add("Link", linkPB)
	w.Header().Set("Content-type", "text/plain")
	fmt.Println(r.URL.Path[1:])
	http.ServeFile(w, r, fmt.Sprintf("./static/%s", r.URL.Path[1:]))
}

// RenderWithProv shows the prov of a resource
// right now it just hist getProvRecord which returns a generic same for all
// record (since I have no prov data stood up now beyond testing stuff)
func RenderWithProv(w http.ResponseWriter, r *http.Request) {
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
	// do something with the POST data
	// likely convert to triples and write to some end point...

	w.WriteHeader(http.StatusNoContent)

	// w.Write([]byte("Thanks for your contribution"))  //  we are 204..  no need for body content
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
