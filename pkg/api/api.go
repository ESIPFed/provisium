package api 

import (
//	"fmt"
	"io/ioutil"
	restful "github.com/emicklei/go-restful"
	util "lab.esipfed.org/provisium/pkg/util"
	kv "lab.esipfed.org/provisium/pkg/kv"

)

// New fires up the services inside textsearch
func New() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/api/v1/api").
		Doc("main pingback service").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// add in start point and length cursors
	service.Route(service.POST("/newKey").To(NewKey).
		Doc("Internal Service for Creating New API Keys").
		Param(service.QueryParameter("key", "new api key")).
		Param(service.QueryParameter("password", "admin password for creating new api keys")).
		Param(service.QueryParameter("name", "name associated with this api key")).
		ReturnsError(400, "Unable to handle request", nil).
		Consumes("text/plain").
		Operation("NewKey"))

	return service
}

// Handles Creating New API Keys 
func NewKey(request *restful.Request, response *restful.Response) {

	// Read the request body and get the parameters
	b,err := ioutil.ReadAll(request.Request.Body)
	util.WriteToLog("info", "New API Key Creation Error: ", err)

	// Get the new API key and password from the request
	key,password,name := util.GetKeyAndPass(string(b))

	// Validate the password 
	if ( password == util.Configuration.ApiKeys ) {

		// write to key value store
		kv.WriteApiKeyAndName(util.Configuration.KvStoreAPI, key, name)

		// send a response
		response.Write( []byte("New API Key Created Successfully") )  
		
	} else { response.Write( []byte("Admin Password Not Valid. No API Key Created.") ) }

}
