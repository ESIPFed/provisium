package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	// TODO move to a proper namespace in earthcube ns
	"lab.esipfed.org/provisium/internal/pingback"
	"lab.esipfed.org/provisium/internal/prov"
	"lab.esipfed.org/provisium/internal/search"
	"lab.esipfed.org/provisium/internal/utils"

	"github.com/gorilla/mux"
	minio "github.com/minio/minio-go"
)

var minioVal, portVal, accessVal, secretVal, bucketVal string
var sslVal bool

func init() {
	akey := os.Getenv("MINIO_ACCESS_KEY")
	skey := os.Getenv("MINIO_SECRET_KEY")

	utils.SPARQLurl = os.Getenv("PROVISIUM_SPARQL")

	flag.StringVar(&minioVal, "address", "localhost", "FQDN for server")
	flag.StringVar(&portVal, "port", "9000", "Port for minio server, default 9000")
	flag.StringVar(&accessVal, "access", akey, "Access Key ID")
	flag.StringVar(&secretVal, "secret", skey, "Secret access key")
	flag.StringVar(&bucketVal, "bucket", "provisium", "The configuration bucket")
	flag.BoolVar(&sslVal, "ssl", false, "Use SSL boolean")
}

func main() {
	// Load configurations
	flag.Parse()
	minioClient := utils.MinioConnection(minioVal, portVal, accessVal, secretVal, bucketVal, sslVal)

	// TODO Do a quick check before moving on
	// utils.checkMinio()
	// utils.checkGraph()

	// Make a new router
	router := mux.NewRouter()
	router.HandleFunc("/prov/doc", prov.HostProvURI).Methods("POST")
	router.HandleFunc("/prov/id/{id}", prov.ID).Methods("GET")
	router.HandleFunc("/prov/doc/{id}", prov.Doc).Methods("GET")
	router.HandleFunc("/prov/id/{id}/search", search.Search).Methods("GET")
	router.Handle("/prov/id/{id}/pingback", minioHandler(minioClient, pingback.PostPing)).Methods("POST")
	router.Handle("/prov/id/{id}/pingback", minioHandler(minioClient, pingback.GetPings)).Methods("GET")
	router.Handle("/prov/id/{id}/pingback/{hash}", minioHandler(minioClient, pingback.GetPingID)).Methods("GET")

	log.Println("Provisium Started :6789")
	log.Fatal(http.ListenAndServe(":6789", router))
}

func minioHandler(minioClient *minio.Client,
	f func(minioClient *minio.Client, w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { f(minioClient, w, r) })
}
