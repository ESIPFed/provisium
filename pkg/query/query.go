package query

import (
	restful "github.com/emicklei/go-restful"
)

// New fires up the services inside textsearch
func New() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/api/v1/query").
		Doc("main pingback service").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// TODO  services needed
	// - query by search term  (return results from pingback, not object store)
	// - query by what else?

	// add in start point and length cursors
	service.Route(service.GET("/search").To(PBSearch).
		Doc("Simple search service").
		Param(service.QueryParameter("q", "query string")).
		ReturnsError(400, "Unable to handle request", nil).
		Operation("PBSearch"))

	return service
}

// PBSearch First test function..   opens each time..  not what we want..
// need to open indexes and maintain state
func PBSearch(request *restful.Request, response *restful.Response) {

	response.Write([]byte("Got it"))
}
