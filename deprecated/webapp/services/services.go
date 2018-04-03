package services

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/knakk/rdf"

	kv "lab.esipfed.org/provisium/webapp/kv"
)

// PingBackCatalog returns prov for a given document ID as a service defined in the W3C Prov-aq NOTE
// This is RDF prov..  not a URI list..
// http://www.example.com/provenance/service?target={+uri}{&steps}
// as noted in the docs.  Optional steps can be, for us, scope [curated, community]
func PingBackCatalog(w http.ResponseWriter, r *http.Request) {
	ID := r.FormValue("target")
	urlElements := strings.Split(ID, "/")

	events, err := kv.GetProvLog(urlElements[len(urlElements)-1])
	if err != nil {
		log.Println("error getting info..  need to 3** a return and log")
	}

	var buffer bytes.Buffer // a buffere for what we are about to collect
	for k, _ := range events {
		content, _, err := kv.GetProvDetails(k)
		if err != nil {
			log.Println("error getting info..  need to 3** a return and log")
		}
		buffer.WriteString(content + "\n")

	}

	fmt.Println(buffer.String()) // make as small local function to return a unique set of strings

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(buffer.String()))
}

// PingBackEvents will get the information we know about a prov ID
// and return it as something like RDF (JSON-LD, turtle..  not sure)
func PingBackEvents(w http.ResponseWriter, r *http.Request) {
	// I really only deal with the UUID, so I need to strip this off.  It would likely
	// be better to do as the NOTE notes (sigh) and regsiter full qualified addresses
	// This is bad form with "logic in code", and I should pass what I expect to get
	ID := r.FormValue("target")
	urlElements := strings.Split(ID, "/")

	events, err := kv.GetProvLog(urlElements[len(urlElements)-1])
	if err != nil {
		log.Println("error getting info..  need to 3** a return and log")
	}

	jsonString, err := json.MarshalIndent(events, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

// PingBackEvents will get the information we know about a prov ID
// and return it as something like RDF (JSON-LD, turtle..  not sure)
func PingBackContent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["ID"]

	content, _, err := kv.GetProvDetails(ID)
	if err != nil {
		log.Println("error getting info..  need to 3** a return and log")
	}

	// I could get the content type..  it is in the KV and returned above, but I ignore it
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(content))
}

// Provenance will get the information we know about a prov ID
// and return it as something like RDF (JSON-LD, turtle..  not sure)
func Provenance(w http.ResponseWriter, r *http.Request) {

	ID := r.FormValue("target")
	fmt.Printf("Want to get prov for %s\n", ID)

	sparqlresults := GetProv(ID)

	log.Println(sparqlresults.Results.Bindings)
	bindings := sparqlresults.Results.Bindings // map[string][]rdf.Term
	binding0 := bindings[0]                    // dumb hack..   we return an array, but only get/want one item..   laugh at me here

	var srmap = map[string]string{}
	srmap["uri"] = ID
	srmap["type"] = binding0["type"].Value
	srmap["qa"] = binding0["qa"].Value
	srmap["at"] = binding0["at"].Value
	srmap["qatype"] = binding0["qatype"].Value
	srmap["attype"] = binding0["attype"].Value
	srmap["agent"] = binding0["agent"].Value
	srmap["role"] = binding0["role"].Value

	log.Println(srmap)

	// 2018/01/08 13:43:11 map[
	// at:t16
	// qatype:http://www.w3.org/ns/prov#attribution
	// attype:http://www.w3.org/ns/prov#attribution
	// agent:http://esipfed.org/org
	// role:http://www.aurole.org/Publisher
	// type:http://www.w3.org/ns/prov#entity
	// qa:t16
	// ]

	// text template
	provtemplate := `@prefix prov: <http://www.w3.org/ns/prov#> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .

<{{.uri}}>
  a <{{.type}}> ;
  prov:qualifiedAttribution _:genid1 ;
  prov:wasAttributedTo _:genid1 .

<http://esipfed.org/org>
  rdf:label "ESIP Lab" ;
  a prov:Agent, prov:org .

_:genid1
  a prov:attribution ;
  prov:agent {{.agent}} ;
  prov:hadRole {{.role}} .
  `

	t := template.Must(template.New("t1").Parse(provtemplate))
	err := t.Execute(w, srmap) // send to a buffer and return the buffer
	if err != nil {
		log.Print(err)
	}
}

// GET RID OF THIS!!! and its dopplerganger in handlers
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
