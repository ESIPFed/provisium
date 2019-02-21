package datastructs

import (
	"encoding/json"

	"lab.esipfed.org/provisium/pkg/kv"
)

// PingBack is where you send back URI-List following the PROV-AQ pattern
// However, since we sotre mimetype, one could take anything like RDF nquads or
// query or SPARQL
type PingBack struct {
	ID       string // UUID for the event
	Mimetype string
	Hash     string
	Body     string
	Date     string
	APIKey   string
}

// ProvObject holds prov records that a agent might generate.  However, they might not
// be able to host it.  So this is where someone could place prov, and get a URL that
// they can use in pingback
type ProvObject struct {
	Hash        string // sha1 string  This is the ID for objects
	Body        string // should be a valid RDF document
	Mimetype    string // needed?   could note seriealization
	Date        string
	APIKey      string
	Depreciated bool
}

func (pb *PingBack) exists() (bool, error) {
	// pb.Hash exists...
	// UUID exist..
	// pb.Body hash///
	return false, nil
}

func (po *ProvObject) exists() (bool, error) {
	// po.Hash exists...
	return false, nil
}

func (pb *PingBack) save(ID string) error {
	// call KV and save this struct in the KV
	// Marshal user data into bytes.
	buf, err := json.Marshal(pb)
	if err != nil {
		return err
	}

	err = kv.StorePB(ID, buf)
	if err != nil {
		return err
	}

	return nil
}