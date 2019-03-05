package pingback

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	minio "github.com/minio/minio-go"
)

// PostPing receives a pingback, validates and archives to an object store
func PostPing(minioClient *minio.Client, w http.ResponseWriter, r *http.Request) {
	log.Println("in post ping")

	// Read the POST
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error on body parameter read %v \n", err)
		log.Println(err)
		w.WriteHeader(422)
		fmt.Fprintf(w, "")
		return
	}

	// TODO move tihs to its own function

	h := sha1.New()
	h.Write([]byte(body))
	bs := h.Sum(nil)
	bss := fmt.Sprintf("%x", bs) // better way to convert bs hex string to string?

	// get the value set for this prov originally
	pa := strings.Split(r.URL.Path, "/")

	// objectName := fmt.Sprintf("%s/%s.jsonld", up.Path, bss)
	objectName := fmt.Sprintf("%s/%s", pa[len(pa)-2], bss) // note getting the SID at index -2 from the URL
	contentType := "application/ld+json"                   // TODO   this is WRONG..  get from the request!!!!
	b := bytes.NewBufferString(string(body))               // TODO need a string buffer from []byte

	usermeta := make(map[string]string) // what do I want to know?
	usermeta["url"] = r.Host
	usermeta["sha1"] = bss
	bucketName := "provisium"

	// Upload the file with FPutObject
	n, err := minioClient.PutObject(bucketName, objectName, b, int64(b.Len()), minio.PutObjectOptions{ContentType: contentType, UserMetadata: usermeta})
	if err != nil {
		log.Printf("%s", objectName)
		log.Println(err)
	}
	log.Printf("Uploaded Bucket:%s File:%s Size %d\n", bucketName, objectName, n)

	// TODO at this point call pingtograph and get back the triples (should be pingtotriples)
	// TODO then call jena loader to insert the triples...

	//Content-Type: application/ld+json; charset=utf-8; profile="http://www.w3.org/ns/json-ld#expanded"
	w.Header().Add("Content-Type", "text/plain; charset=utf-8; ")
	msg := fmt.Sprintf("http://localhost:6789/prov/id/%s/pingback/%s", pa[len(pa)-2], bss)
	fmt.Fprintf(w, "%s", msg)
}
