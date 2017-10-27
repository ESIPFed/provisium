package kv

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// newProvEvent must address a range of actions.  On a new event
// we need to record the CLF of the event, the prov graph fragment
// and associate the new prov event ID with the document ID
func newProvEvent(docID, provFrag string) error {

	fmt.Printf("For doc %s I am recording a new event", docID)

	provID := "make UUID here"

	// Need to try and make this transactional at some point...
	// Out of scope initially for the project...
	// would likely have to use some roll back on a not nil event
	err := provLog(provID)
	if err != nil {
		log.Printf("ERROR:  could not log the event")
	}
	err = provRecord(provID, provFrag)
	if err != nil {
		log.Printf("ERROR:  could not record the provenance graph fragment")
	}
	err = docToProvID(docID, provID)
	if err != nil {
		log.Printf("ERROR:  could not associate provID to a docID")
	}

	return nil
}

// provLog will record a CLF log event in a KV store
func provLog(provID string) error {
	fmt.Printf("I will log an event %s\n", provID)
	fmt.Printf("can I get the CLF string from gorilla..  or must I make it?\n")

	return nil
}

// provRecord will take a proposed prov event graph fragment and
// first validate it as RDF and against any specified criteria.  It will
// log it to the KV store and send it to the master graph
func provRecord(provID, provFrag string) error {

	// err := validateAsRDF(provFrag)
	// err = storeToDefaultGraph(provFrag)

	fmt.Printf("I will store a prov graph fragment %s\n", provFrag)

	return nil

}

// docToProvID will associagte a prov event ID with a document ID
func docToProvID(docID, provID string) error {
	fmt.Printf("I will associate a docID %s with a provID  %s\n", docID, provID)

	return nil

}

// GetDocIDs get all the files in our holding
func GetDocIDs() []string {
	db := getKVStore()

	var IDs []string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MetaDataBucket"))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			// log.Printf("key=%s, value=%s\n", k, v)
			IDs = append(IDs, string(k))
		}
		return nil
	})

	if err != nil {
		log.Println("Error reading file info from the KV store index.db")
		log.Println(err)
	}

	err = db.Close()
	if err != nil {
		log.Println("Error closing database index.db")
		log.Println(err)
	}

	return IDs
}

// GetResMetaData will get the metadata for a dataset
func GetResMetaData(docID string) (string, error) {
	fmt.Printf("I will get the info for docID %s \n", docID)
	db := getKVStore()

	var jsonld string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MetaDataBucket"))
		v := b.Get([]byte(docID))
		jsonld = string(v)
		return nil
	})

	if err != nil {
		log.Println("Error reading file info from the KV store index.db")
	}

	err = db.Close()
	if err != nil {
		log.Println("Error closing database index.db")
		log.Println(err)
	}

	return jsonld, err
}

func getKVStore() *bolt.DB {
	db, err := bolt.Open("./kvStores/index.db", 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	return db
}
