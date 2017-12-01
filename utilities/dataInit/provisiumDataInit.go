package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/kazarena/json-gold/ld"
)

type fileInfo struct {
	Name string
	UUID string
}

// FramedFileInfo holds the results of our framing function
// the _ naming violates Go coding patterns
type FramedFileInfo struct {
	_id        string `json:"@id"`
	_type      string `json:"@type"`
	Identifier struct {
		_id   string `json:"@id"`
		_type string `json:"@type"`
		Value string `json:"value"`
	} `json:"identifier"`
	Name string `json:"name"`
}

func main() {
	// Open the my.db data file in your current directory.
	// It will be created if it doesn't exist.
	db, err := bolt.Open("index.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = initBuckets(db)
	if err != nil {
		log.Println(err)
		panic(err) // might as well..  things have gone bad...
	}

	// files, err := getMetaDataFile(".")
	files, err := filePathGlob(".")
	if err != nil {
		log.Println(err)
	}

	// Frame the JSON-LD and extract the
	// UUID and filename to place in the KV store
	// array of fileInfo strucres
	// fileInfo = []fileInfo{}
	for k := range files {
		dat, err := ioutil.ReadFile(files[k])
		if err != nil {
			fmt.Printf("Error reading file %s\n", files[k])
		}
		fi := frameDoc(dat)
		fmt.Printf("Found %s\n", files[k])
		fmt.Println(fi[0].Identifier.Value)
		fmt.Println(fi[0].Name)

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("FileBucket"))
			err := b.Put([]byte(fi[0].Identifier.Value), []byte(fi[0].Name))
			return err
		})

		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MetaDataBucket"))
			err := b.Put([]byte(fi[0].Identifier.Value), dat)
			return err
		})

	}
}

// frameDoc using JSON-LD frame API to generate a JSON we can easily
// map to our data struct
func frameDoc(dat []byte) []FramedFileInfo {
	proc := ld.NewJsonLdProcessor()
	options := ld.NewJsonLdOptions("")

	frame := map[string]interface{}{
		"@context": map[string]interface{}{
			"@vocab": "http://schema.org/",
		},
		"@type":     "Dataset",
		"@explicit": true,
		"name":      map[string]interface{}{},
		"identifier": map[string]interface{}{
			"@explicit": true,
			"@type":     "PropertyValue",
			"value":     map[string]interface{}{},
		},
	}

	var myInterface interface{}
	err := json.Unmarshal(dat, &myInterface)
	if err != nil {
		log.Println("Error when transforming JSON-LD document to interface:", err)
	}

	framedDoc, err := proc.Frame(myInterface, frame, options) // do I need the options set in order to avoid the large context that seems to be generated?
	if err != nil {
		log.Println("Error when trying to frame document", err)
	}

	graph := framedDoc["@graph"]
	jsonm, err := json.MarshalIndent(graph, "", " ")
	if err != nil {
		log.Println("Error trying to marshal data", err)
	}

	dss := make([]FramedFileInfo, 0)
	err = json.Unmarshal(jsonm, &dss)
	if err != nil {
		log.Println("Error trying to unmarshal data to struct", err)
	}

	// log.Printf("This is the json now:\n  %v\n", string(jsonm))
	// log.Print(dss[0].Identifier.Value)
	// log.Print(dss[0].Name)

	// return fileInfo{Name: "filename", UUID: "UUID value"}

	return dss
}

func filePathGlob(directory string) ([]string, error) {
	pattern := fmt.Sprintf("%s/*.json", directory)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println(err)
	}
	return matches, err
}

func initBuckets(db *bolt.DB) error {

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("FileBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("MetaDataBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	return err
}
