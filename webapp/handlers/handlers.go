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
	"github.com/tomnomnom/linkheader"
	kv "lab.esipfed.org/provisium/webapp/kv"
)

type PageData struct {
	SchemaOrg string
	EventLog  map[string]string
	ProvRDF   string
	Host      string
	TargetURI string
	UUID      string
	DataFile  string
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
	events, err := kv.GetProvLog(ID) // TODO, this should return an error too
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	targeturi := r.Host + r.URL.String()
	pd := PageData{SchemaOrg: so, EventLog: events, Host: r.Host, TargetURI: targeturi, UUID: ID, DataFile: getDataFileName(ID)} // struct to pass to the page

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
	w.Write([]byte(getProvRecord())) // TODO.. make this return something sorta real pulled from a KV store or something....
}

// ProvPingback Handles the PROV pingback on a resource, ref: https://www.w3.org/TR/prov-aq/ Section 5
func ProvPingback(w http.ResponseWriter, r *http.Request) {

	// If we are not POST..  get out
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// require uri list for pingback...  otherwise, get out
	contentType := r.Header.Get("Content-Type")
	if contentType != "text/uri-list" {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	// read the body..   if trouble..   get out
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var status int
	var docerr error

	// The pingback can have link header so look for/get them. ref: Section 5, examples 13 and 14 of https://www.w3.org/TR/prov-aq/
	links := r.Header.Get("Link") // NOTE  via https://tools.ietf.org/html/rfc5988#section-5.5 multiple entries exposed by comma
	if links != "" {
		fmt.Println(links)
		linksitems := linkheader.Parse(links)
		for _, link := range linksitems {
			// validate the links..  store them in a KV store designed for these relations
			// _, err := url.ParseRequestURI(link.URL) // validate this is a URL
			// TODO if no anchor...  set to target URI
			status, docerr = recordLinkItem(link.URL, link.Rel, r.RemoteAddr, link.Param("anchor"))
		}
	}

	// check the length of the body..  if set to 0 or evaluates to 0, don't bother to record anything (duh...)
	contentLength := r.Header.Get("Content-Length") // a string, not an int
	evalBodyLength := len(body)                     // an int, not a string

	if contentLength == "0" || evalBodyLength == 0 {
		fmt.Println("Body set to or is 0..  do nothing")
	} else {
		// record the body
		status, docerr = recordBody(body, contentType, r)
	}

	fmt.Println(docerr)

	// We REALLY made it..   tell them all is OK with 204 empty
	w.WriteHeader(status)
}

func recordLinkItem(url, rel, ip, anchor string) (int, error) {

	fmt.Printf("URL: %s; Rel: %s Anchor: %s\n", url, rel, anchor)

	err := kv.NewProvEvent(anchor, url, ip, rel)
	if err != nil {
		fmt.Println("Error trying to record the uploaded prov")
		// w.WriteHeader(http.StatusUnprocessableEntity)
		return http.StatusUnprocessableEntity, err
	}

	return http.StatusNoContent, nil
}

func recordBody(body []byte, contentType string, r *http.Request) (int, error) {

	scanner := bufio.NewScanner(strings.NewReader(string(body)))
	var URLError error
	for scanner.Scan() {
		_, err := url.ParseRequestURI(scanner.Text()) // validate this is a URL
		if err != nil {
			URLError = err
		}
		fmt.Printf("URL %s report: %v\n", scanner.Text(), err)
	}
	if URLError != nil {
		log.Printf("Error in the sent URLs: %v", URLError)
		// w.WriteHeader(http.StatusUnprocessableEntity)
		return http.StatusUnprocessableEntity, URLError
	}

	// We made it, record is valid..  try and record it now
	fmt.Printf("Recording a new prov event for %s\n", r.URL.Path[1:])
	pathElements := strings.Split(r.URL.Path[1:], "/")
	docID := pathElements[2]

	err := kv.NewProvEvent(docID, string(body), r.RemoteAddr, contentType)
	if err != nil {
		fmt.Println("Error trying to record the uploaded prov")
		// w.WriteHeader(http.StatusUnprocessableEntity)
		return http.StatusUnprocessableEntity, err
	}

	return http.StatusNoContent, nil

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

// GetData will get the data file
func getDataFileName(UID string) string {

	content, err := kv.GetResData(UID)
	if err != nil {
		log.Println("error getting file content..  need to 3** a return and log")
	}

	return content
}
