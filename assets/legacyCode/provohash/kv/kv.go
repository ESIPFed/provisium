package kv

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/twinj/uuid"
)

// NewProvHash  get a hash for a prov, but checks to ensure it's not there already
func NewProvHash(provFrag, remoteAddress, contentType string) (string, error) {
	// DEPRECATED approach..   issues with hashing and event provenance (ironically)
	// provHash := getMD5Hash(provFrag)

	provUUID := getUUID()

	// This did check to see if a key by hash existed..
	//  The UUID will obviously not..   but we could check the body hash ...
	// need to refactor the code a bit to do that and not sure what the use case would
	// be to do that....
	if checkForKey(provUUID) {
		log.Print("PROV already registered")
		return "", nil
	}
	log.Printf("checked for prov %s", provUUID)

	db := getKVStoreRW()

	// TODO..  connect these three updates into a single transaction wrapper
	// Log the event
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("LogBucket"))
		logEvent := fmt.Sprintf("%s, %s, %s", remoteAddress, time.Now().String(), contentType)
		err := b.Put([]byte(provUUID), []byte(logEvent))
		return err
	})

	// Record the Prov
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ProvBucket"))
		err := b.Put([]byte(provUUID), []byte(provFrag))
		return err
	})

	db.Close()

	return provUUID, err
}

// Not really doing anything now that the key is not an MD5..   but we
// still could check for identical prov entry...  in cased that is of interest to know...
func checkForKey(keyToCheck string) bool {
	db := getKVStoreRO()
	keyExists := false

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ProvBucket"))
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("compare %s to %s, \n", keyToCheck, string(k))
			if keyToCheck == string(k) {
				keyExists = true
			}
			return nil
		})
		return nil
	})

	db.Close()

	return keyExists
}

//GetProvDetails Return the entry for a specific prov record
func GetProvDetails(provID string) (string, string, error) {

	fmt.Printf("Request URI list entries for provID %s \n", provID)
	db := getKVStoreRO()

	var provEntry string
	var contentType string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ProvBucket"))
		b2 := tx.Bucket([]byte("LogBucket"))

		v := b.Get([]byte(provID))
		provEntry = string(v)

		v2 := b2.Get([]byte(provID))
		logLine := strings.Split(string(v2), ",") // 3rd entry is content-type, see NewProvEvent write to this bucket
		if len(logLine) == 3 {
			contentType = logLine[2]
		} else {
			contentType = "text/plain"
		}

		log.Println(logLine)

		return nil
	})

	if err != nil {
		log.Println("Error reading from Buckets")
	}

	err = db.Close()
	if err != nil {
		log.Println("Error closing database index.db")
		log.Println(err)
	}

	return provEntry, contentType, nil
}

func GetListing() []string {
	db := getKVStoreRO()

	var IDs []string
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ProvBucket"))
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

func getUUID() string {
	uuid := uuid.NewV4()
	return uuid.String()
}

func getMD5Hash(text string) string {
	h := md5.New()
	io.WriteString(h, text)
	return fmt.Sprintf("%x", h.Sum(nil))
}
