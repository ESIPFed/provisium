package util

import (
	"os"
	"time"
	"strings"
	"io/ioutil"
	"os/exec"
	"crypto/sha256"
	"encoding/base64"
	"github.com/tkanos/gonfig"
	log "github.com/sirupsen/logrus"
	structs "lab.esipfed.org/provisium/pkg/datastructs"
)

// Helper function for logging
func WriteToLog(logLevel string, message string, err error) {

	if err != nil {
		if logLevel == "fatal" {
			log.Fatal(message, err)
		} else {
			log.Info(message, err)
        	}
	}

}

// Helper function to read the configuration 
var Configuration = structs.Configuration{}
func ReadConfig() (error) {
	//configuration := Configuration{}
	err := gonfig.GetConf("./config/config.json", &Configuration)
	return err
}

// Helper function to loop over a directory and 
// insert all PROV docs found into the triple store
func RunScript(cmd string, args []string) error {

	err := exec.Command(cmd,args...).Run()
	WriteToLog("fatal", "Error executing script: " + cmd, err)

	return err

}

// Helper function get PROV file extension
func GetExtension(encoding string) string {

	ext := ""
	if encoding == "turtle" {
		ext = ".ttl"
	} else if encoding == "rdf-xml" {
		ext = ".rdf"
	} else if encoding == " n-triples" {
		ext = ".nt"
	} else if encoding == "n-quads" {
		ext = ".nq"
	}

	return ext
}

// Helper function to write PROV to a file
func WriteToFile(dir, prov, encoding, apiKey string) (string, string, error) {

	ext := " "
	if encoding == "turtle" {
		ext = ".ttl"
	} else if encoding == "rdf-xml" {
		ext = ".rdf"
	} else if encoding == "n-triples" {
		ext = ".nt"
	} else if encoding == "n-quads" {
		ext = ".nq"
	}

	// get an SHA 256 hash of the time to use as a unique filename
	h := sha256.New()
	data := []byte(prov)
	h.Write( []byte(time.Now().String()) )
	sha := base64.URLEncoding.EncodeToString( h.Sum(nil) )

	// does the local directory exist? If not, create it
	localDir := dir + apiKey + "/"
	if _, err := os.Stat(localDir); os.IsNotExist(err) { os.Mkdir(localDir, 0777) }

	// write the file to the local directory 
	file := localDir + sha + ext 
	web := Configuration.URL + sha + ext
	err := ioutil.WriteFile(file, data, 0644)
	return web, file, err	
}

// Helper function to parse API request
// for creating new API keys
func GetKeyAndPass(s string) (string, string, string) {

	var key, psswd, name string = "", "", ""
	parts := strings.Split(s,"&")
	v1 := strings.Split(parts[0],"=")
	v2 := strings.Split(parts[1],"=")
	v3 := strings.Split(parts[2],"=")
	
	if v1[0] == "key" {
	    key = v1[1]
	} else if v2[0] == "key" {
	    key = v2[1]
	} else { key = v3[1] }

	if v1[0] == "password" {
	    psswd = v1[1]
	} else if v2[0] == "password" {
	    psswd = v2[1]
	} else { psswd = v3[1] }

	if v1[0] == "name" {
	   name = v1[1]
	} else if v2[0] == "name" {
	   name = v2[1]
	} else { name = v3[1] }

	return key, psswd, name

}

// Helper function to parse API POST input
// and extract API Key and PROV text
func GetKeyAndProv(s string) (string, string, string) {

	var key, prov, encoding string = "", "", ""
	parts := strings.Split(s,"&")
	if ( len(parts) == 3 ) {
		v1 := strings.Split(parts[0],"=")
		v2 := strings.Split(parts[1],"=")
		v3 := strings.Split(parts[2],"=")

		if v1[0] == "key" {
			key = v1[1]
		} else if v2[0] == "key" {
			key = v2[1]
		} else { key = v3[1] }

		if v1[0] == "prov" {
			prov = v1[1]
		} else if v2[0] == "prov" {
			prov = v2[1]
		} else { prov = v3[1] }

		if v1[0] == "encoding" {
			encoding = v1[1]
		} else if v2[0] == "encoding" {
			encoding = v2[1]
		} else { encoding = v3[1] }
	}

	return key, prov, encoding 
}

// Helper function to parse API POST input
// and extract API Key, PROV, encoding, and url
func GetKeyAndUrl(s string) (string, string, string, string) {

	var key, prov, encoding, url string = "", "", "", ""
	parts := strings.Split(s,"&")
	if ( len(parts) == 4 ) {
		v1 := strings.Split(parts[0],"=")
		v2 := strings.Split(parts[1],"=")
		v3 := strings.Split(parts[2],"=")
		v4 := strings.Split(parts[3],"=")

		if v1[0] == "key" {
			key = v1[1]
		} else if v2[0] == "key" {
			key = v2[1]
		} else if v3[0] == "key" {
			key = v3[1]
		} else { key = v4[1] }

		if v1[0] == "prov" {
			prov = v1[1]
		} else if v2[0] == "prov" {
			prov = v2[1]
		} else if v3[0] == "prov" {
			prov = v3[1]
		} else { prov = v4[1] }

		if v1[0] == "encoding" {
			encoding = v1[1]
		} else if v2[0] == "encoding" {
			encoding = v2[1]
		} else if v3[0] == "encoding" {
			encoding = v3[1]
		} else { encoding = v4[1] }

		if v1[0] == "url" {
			url = v1[1]
		} else if v2[0] == "url" {
			url = v2[1]
		} else if v3[0] == "url" {
			url = v3[0]
		} else { url = v4[1] }
	}

	return key, prov, encoding, url
}
