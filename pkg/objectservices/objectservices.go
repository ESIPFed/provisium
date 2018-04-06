package objectservices

import (
	"io/ioutil"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

// New fires up the services inside textsearch
func New() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/api/v1/objects").
		Doc("main pingback service").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// TODO  services needed
	// - store prov (checks that it is RDF)
	// - get prov by HASH
	// - get prov entries  (start, step)
	// - query .. return results from object store only..  not pingback
	// - delete ?  or perhaps not...  (set a deprecated flag?)

	// add in start point and length cursors
	service.Route(service.POST("/post").To(PostObject).
		Doc("Base W3C PROV-AQ pingback implementation").
		Param(service.BodyParameter("body", "The body containing a list of URIs")).
		Consumes("text/uri-list").
		Produces("text/plain").
		ReturnsError(400, "Unable to handle request", nil).
		Operation("PutObject"))

	return service
}

// PostObject First test function..   opens each time..  not what we want..
// need to open indexes and maintain state
func PostObject(request *restful.Request, response *restful.Response) {
	// read the body..   if trouble..   get out
	b, err := ioutil.ReadAll(request.Request.Body) // here we need to get to the raw REQUEST to read the raw BODY
	if err != nil {
		log.Printf("Error reading body: %v", err)
		response.WriteError(400, err)
		return
	}

	hl := request.HeaderParameter("Link")
	cl := request.HeaderParameter("Content-Length")
	bl := len(b)

	log.Printf("Header Link: %s", hl)
	log.Printf("Content Length: %s", cl)
	log.Printf("Body Length: %d", bl)
	log.Printf("URI list: %s", (string(b)))

	if cl == "0" || bl == 0 {
		log.Println("During testing we will not exit on content or body length 0 status")
	}

	// todo: write the URI list to KV store since we made it this far
	response.Write([]byte("Got it"))
}
