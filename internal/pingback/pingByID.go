package pingback

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	minio "github.com/minio/minio-go"
)

// GetPingID collects all the pings and returns them in a human format
func GetPingID(minioClient *minio.Client, w http.ResponseWriter, r *http.Request) {
	// minioClient := utils.MinioConnection()

	pa := strings.Split(r.URL.Path, "/")
	object, err := minioClient.GetObject("provisium", fmt.Sprintf("%s/%s", pa[3], pa[5]), minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8; profile=\"http://www.w3.org/ns/json-ld#expanded\"")

	n, err := io.Copy(w, object)
	if err != nil {
		log.Println("Issue with writing file to http response")
		log.Println(err)
	}
	log.Printf("Wrote file with bytes %d\n", n)
}
