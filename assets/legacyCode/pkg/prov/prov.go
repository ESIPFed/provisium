package prov 

import (
//	"fmt"
	"os"
	"strings"
	"io/ioutil"
	restful "github.com/emicklei/go-restful"
	util "lab.esipfed.org/provisium/pkg/util"
	kv "lab.esipfed.org/provisium/pkg/kv"

)

// New fires up the services inside textsearch
func New() *restful.WebService {
	service := new(restful.WebService)
	service.Path("/api/v1/prov").
		Doc("provenance submission service").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// add in start point and length cursors
	service.Route(service.POST("/submission").To(Submission).
		Doc("Submission Service for PROV").
		Param(service.QueryParameter("prov", "provenance text")).
		Param(service.QueryParameter("encoding", "provenance encoding (turtle, rdf-xml, n-triple, or n-quad)")).
		Param(service.QueryParameter("key", "API key")).
		ReturnsError(400, "Unable to handle request", nil).
		Consumes("text/plain").
		Operation("Submission"))

	// add in start point and length cursors
	service.Route(service.POST("/update").To(Update).
		Doc("Update previously submitted PROV").
		Param(service.QueryParameter("prov", "provenance text")).
		Param(service.QueryParameter("encoding", "provenance encoding (turtle, rdf-xml, n-triple, or n-quad)")).
		Param(service.QueryParameter("key", "API key")).
		Param(service.QueryParameter("url", "url of the file you want to update. this was returned upon original submission")).
		ReturnsError(400, "Unable to handle request", nil).
		Consumes("text/plain").
		Operation("Update"))

	return service
}

// Handles PROV submissions 
func Submission(request *restful.Request, response *restful.Response) {

	// Read the request body and get the parameters
	b,err := ioutil.ReadAll(request.Request.Body)
	util.WriteToLog("info", "PROV submission error: ", err)

	// Get the API key and PROV from the request
	key,prov,encoding := util.GetKeyAndProv(string(b))

	// Did we get all of the required parameters?
	inputError := false
	if ( key == "" || prov == "" || encoding == "" ) { inputError = true }
	if ( !inputError ) {
	
		// Validate the API key
		keyIsValid := kv.ValidateApiKey(key, util.Configuration.KvStoreAPI)
		if keyIsValid {

			// Make sure this is an encoding we understand
			encodingError := false
			l := strings.ToLower(encoding)
			if ( l != "turtle" && l != "rdf-xml" && l != "n-quads" && l != "n-triples" ) { encodingError = true }

			// Validate the PROV document

			// If no error then write to file
			if ( !encodingError ) {

				// Write PROV to file
				url, file, err := util.WriteToFile(util.Configuration.WebStorageDir, prov, encoding, key)
				util.WriteToLog("fatal", "PROV write error: ", err)
				response.Write( []byte(url) )

				// Write to key/value store - url, api key
				err = kv.WriteUrlAndKey(util.Configuration.KvStoreAPI, url, key)
				util.WriteToLog("fatal", "Error writing url/key pair to KV store: ", err)

				// Insert the PROV into the triple store
				cmd := "exec/vload"
				ext := util.GetExtension(encoding) 
				graphUri := util.Configuration.GraphUriBase + key 
				args := []string{ ext[1:], file, graphUri}
				err = util.RunScript(cmd, args)
				util.WriteToLog("fatal", "Error inserting into triple store: ", err)

			} else { response.Write( []byte("Encoding Not Understood. Valid Encodings {turtle, rdf-xml, n-quads, n-triples}") ) }
		
		} else { response.Write( []byte("API Key Not Valid. Submission Not Accepted")) }

	} else { response.Write( []byte("API Usage Error: 3 parameters required (apiKey, prov, encoding)") ) } 

}

// Handles PROV Updates
func Update(request *restful.Request, response *restful.Response) {

	// Read the request body and get the parameters
	b,err := ioutil.ReadAll(request.Request.Body)
	util.WriteToLog("info", "PROV update error: ", err)

	// Get the API key, PROV, encoding, and URL from the request
	key,prov,encoding,url := util.GetKeyAndUrl(string(b))

	// Did we get all of the required API parameters?
	inputError := false
	if ( key == "" || prov == "" || encoding == "" || url == "" ) { inputError = true }
	if ( !inputError ) {

		// Validate the API key
		keyIsValid := kv.ValidateApiKey(key, util.Configuration.KvStoreAPI)
		if keyIsValid {

			// Make sure this is an encoding we understand
			encodingError := false
			l := strings.ToLower(encoding)
			if ( l != "turtle" && l != "rdf-xml" && l != "n-quads" && l != "n-triples" ) { encodingError = true }

			// Validate the PROV document

			// if no error then delete the current file and replace with the new one
			if ( !encodingError ) {

				// Delete old file
				parts := strings.Split(url, "/")
				file := util.Configuration.WebStorageDir + key + "/" + parts[len(parts)-1] + util.GetExtension(encoding)
				err := os.Remove(file)
				util.WriteToLog("fatal", "PROV update - unable to delete file: ", err)

				// Delete old url,api key from key value store
				err = kv.DeleteUrlAndKey(util.Configuration.KvStoreAPI, url)

				// Write PROV to file
				url, file, err := util.WriteToFile(util.Configuration.WebStorageDir, prov, encoding, key)
				util.WriteToLog("fatal", "PROV write error (in Update): ", err)
				response.Write( []byte(url) )

				// Write new url,api key to key value store
				err = kv.WriteUrlAndKey(util.Configuration.KvStoreAPI, url, key)

				// Delete everything from this graph in the triple store
				cmd := "exec/vdelete"
				args := []string{ util.Configuration.GraphUriBase + key }
				err = util.RunScript(cmd, args)
				util.WriteToLog("fatal", "Error deleting from triple store: ", err)

				// Insert all PROV files for this API key
				files, err := ioutil.ReadDir(util.Configuration.WebStorageDir + key + "/")
				util.WriteToLog("fatal", "Error inserting PROV to triple store: ", err)
		    		for _, f := range files {
		    			// get the encoding
		    			parts = strings.Split(f.Name(), ".")
        	    			ext := parts[len(parts)-1]
        	    			args := []string{ ext, file, util.Configuration.GraphUriBase + key }
        	    			err = util.RunScript("exec/vload", args)
					util.WriteToLog("fatal", "Error deleting from triple store: ", err)
    				}

			} else { response.Write( []byte("Encoding Not Understood. Valid Encodings {turtle, rdf-xml, n-quads, n-triples}") ) }

		} else { response.Write([]byte("API Key Not Valid. Submission Not Accepted")) }

	} else { response.Write([]byte("API Usage Error: 4 parameters required (apiKey, prov, encoding, url)")) }

}
