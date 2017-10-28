package kv

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/google/uuid"
)

// newProvEvent must address a range of actions.  On a new event
// we need to record the CLF of the event, the prov graph fragment
// and associate the new prov event ID with the document ID
func NewProvEvent(docID, provFrag string) error {

	provID := uuid.New().String()
	fmt.Printf("For doc %s I am recording a new event %s \n", docID, provID)

	// Need to try and make this transactional at some point...
	// Out of scope initially for the project...
	// would likely have to use some roll back on a not nil event

	db := getKVStoreRW()

	// TODO..  connect these three updates into a single transaction wrapper
	// Log the event
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("LogBucket"))
		logEvent := fmt.Sprintf("Pingback at: %s", time.Now().String())
		err := b.Put([]byte(provID), []byte(logEvent))
		return err
	})

	// Record the Prov
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ProvBucket"))
		err := b.Put([]byte(provID), []byte(provFrag))
		return err
	})

	// Associate DocID and ProvID
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("IDLinkBucket"))
		err := b.Put([]byte(provID), []byte(docID))
		return err
	})

	db.Close()

	return nil
}

func GetProvDetails(provID string) string {
	// TODO..  get all the details about a prov event...
	// TODO..  how to address some events bying type text/uri-list and others
	// being something like turtle or other RDF encodings.
}

// GetProvLog gets all the logged events for a given docID
func GetProvLog(docID string) map[string]string {
	db := getKVStoreRO()

	eventmap := make(map[string]string)

	// Logic needed
	// 1) loop over IDLinkBucket to find all provID that match a value of docID
	// 2) for each provID, pull event (value) from LogBucket
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("IDLinkBucket"))
		b2 := tx.Bucket([]byte("LogBucket"))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if strings.Contains(string(v), docID) {
				v2 := b2.Get(k)
				eventmap[string(k)] = string(v2)
			}
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

	return eventmap

}

// GetDocIDs get all the files in our holding
func GetDocIDs() []string {
	db := getKVStoreRO()

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
	db := getKVStoreRO()

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

func getKVStoreRW() *bolt.DB {
	db, err := bolt.Open("./kvStores/index.db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	return db
}

func getKVStoreRO() *bolt.DB {
	db, err := bolt.Open("./kvStores/index.db", 0666, &bolt.Options{ReadOnly: true})
	if err != nil {
		log.Fatal(err)
	}
	// defer db.Close()

	return db
}

// Init the KV store in case we are starting empty and need some buckets made
// Call from the main program at run time...
func InitKV() error {

	db := getKVStoreRW()

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

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("IDLinkBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("ProvBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("LogBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	db.Close()

	return err

}
