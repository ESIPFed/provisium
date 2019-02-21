package prov

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// UpdateCall test out updates to Jena
func UpdateCall(s []byte) ([]byte, error) {
	url := "http://clear.local:3030/provisium/update"
	// fmt.Println("URL:>", url)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(s))
	req.Header.Set("Content-Type", "application/sparql-update")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return body, err
}
