package pingback

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	minio "github.com/minio/minio-go"
)

// GetPings collects all the pings and returns them in a human format
func GetPings(minioClient *minio.Client, w http.ResponseWriter, r *http.Request) {
	log.Println("get pings")

	// // Create a done channel to control 'ListObjectsV2' go routine.
	doneCh := make(chan struct{})
	defer close(doneCh) // Indicate to our routine to exit cleanly upon return.
	pa := strings.Split(r.URL.Path, "/")

	oa := []minio.ObjectInfo{}

	isRecursive := true
	objectCh := minioClient.ListObjectsV2("provisium", pa[len(pa)-2], isRecursive, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			return
		}
		oa = append(oa, object)
	}

	//Content-Type: application/ld+json; charset=utf-8; profile="http://www.w3.org/ns/json-ld#expanded"
	w.Header().Add("Content-Type", "application/ld+json; charset=utf-8; profile=\"http://www.w3.org/ns/json-ld#expanded\"")
	jld, err := json.MarshalIndent(oa, "", " ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprintf(w, "%s", jld)
}
