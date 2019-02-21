package main

import (
	"net/http"
	"os"
	log "github.com/sirupsen/logrus"
	"lab.esipfed.org/provisium/pkg/kv"
	"lab.esipfed.org/provisium/pkg/prov"
	"lab.esipfed.org/provisium/pkg/api"
	util "lab.esipfed.org/provisium/pkg/util"

	restful "github.com/emicklei/go-restful"
	swagger "github.com/emicklei/go-restful-swagger12"

)

func init() {
	log.SetFormatter(&log.JSONFormatter{}) // Log as JSON instead of the default ASCII formatter.
	log.SetOutput(os.Stdout)               // I override this and set output to file (io.Writer) in main
	log.SetLevel(log.DebugLevel)           // Will log anything that is info or above (debug, info, warn, error, fatal, panic). Default.
}

func main() {

	// Read the provisium config file
	err := util.ReadConfig()
	util.WriteToLog("fatal", "main.go, Unable to read provisium congig file: ", err) 

	// Set up our log file 
	f, err := os.OpenFile(util.Configuration.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	util.WriteToLog("fatal", "main.go, Error opening log file: ", err)
	defer f.Close()
	log.SetOutput(f)

	// set up the KV store
	err = kv.InitKV(util.Configuration.KvStoreAPI)
	util.WriteToLog("fatal", "main.go, Error setting up KV store: ", err)

	wsContainer := restful.NewContainer()

	// CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)

	// Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)

	// Add the services
	wsContainer.Add(prov.New())          
	wsContainer.Add(api.New())

	// Swagger
	swaggerConfig := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		ApiPath:        "/apidocs.json",
		WebServicesUrl: "http://provisium.io"} // localhost
	swagger.RegisterSwaggerService(swaggerConfig, wsContainer)

	// Start up
	log.Printf("Services on localhost:" + util.Configuration.Port)
	server := &http.Server{Addr: ":" + util.Configuration.Port, Handler: wsContainer}
	server.ListenAndServe()
}
