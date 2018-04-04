package datastructs

type PingBack struct {
	ID       string // UUID for the event
	Mimetype string
	Hash     string
	Body     string
	Date     string
}

type ProvObject struct {
	Hash      string // sha1 string
	Body      string // should be a valid RDF document
	Mimetype  string // needed?   could note seriealization
	EventInfo string
	Date      string
}

func (pb *PingBack) exists() (bool, error) {
	// pb.Hash exists...
	return false, nil
}

func (po *ProvObject) exists() (bool, error) {
	// po.Hash exists...
	return false, nil
}

func (pb *PingBack) save() error {
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
