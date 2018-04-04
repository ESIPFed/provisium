package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"lab.esipfed.org/provisium/pkg/kv"
	"lab.esipfed.org/provisium/pkg/objectservices"
	"lab.esipfed.org/provisium/pkg/pingback"
	"lab.esipfed.org/provisium/pkg/query"

	restful "github.com/emicklei/go-restful"
	swagger "github.com/emicklei/go-restful-swagger12"
)

func init() {

	log.SetFormatter(&log.JSONFormatter{}) // Log as JSON instead of the default ASCII formatter.
	log.SetOutput(os.Stdout)               // I override this and set output to file (io.Writer) in main
	log.SetLevel(log.DebugLevel)           // Will log anything that is info or above (debug, info, warn, error, fatal, panic). Default.
}

func main() {
	// Set up our log file for runs...
	f, err := os.OpenFile("./logs/serviceslog.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// set up the KV store
	err = kv.InitKV()
	if err != nil {
		log.Fatal(err) // fatal..  no bucket..  no reason
	}

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
	wsContainer.Add(pingback.New())       // text search services
	wsContainer.Add(query.New())          // text search services
	wsContainer.Add(objectservices.New()) // text search services

	// Swagger
	config := swagger.Config{
		WebServices:    wsContainer.RegisteredWebServices(), // you control what services are visible
		ApiPath:        "/apidocs.json",
		WebServicesUrl: "http://provisium.io"} // localhost:6789
	swagger.RegisterSwaggerService(config, wsContainer)

	// Start up
	log.Printf("Services on localhost:6789")
	server := &http.Server{Addr: ":6789", Handler: wsContainer}
	server.ListenAndServe()
}
