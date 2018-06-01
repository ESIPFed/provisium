package kv

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

// Write api key and name
func WriteApiKeyAndName (kvStore, apiKey, name string) error {

	db := getKVStoreRW(kvStore)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ApiKeys"))
		err := b.Put([]byte(apiKey), []byte(name))
		return err
	})

	db.Close()

	return nil
}

// Delete url and api key from the key store
func DeleteUrlAndKey (kvStore, url string) error {

	db := getKVStoreRW(kvStore)

	db.Update(func(tx * bolt.Tx) error {
		b := tx.Bucket([]byte("FileBucket"))
		err := b.Delete([]byte(url))
		return err
	})

	db.Close()

	return nil

}

// Write url and api key to the key store
func WriteUrlAndKey (kvStore, url, apiKey string) error {

	db := getKVStoreRW(kvStore)

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("FileBucket"))
		err := b.Put([]byte(url), []byte(apiKey))
		return err
	})

	db.Close()

	return nil
}

// ValidateApiKey will validate an API Key (true/false) 
func ValidateApiKey(apiKey string, kvStore string) bool {

	db := getKVStoreRO(kvStore)

	isValid := false
	name := ""
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("ApiKeys"))
		v := b.Get([]byte(apiKey))
		name = string(v)
		return nil
	})

	if err != nil {
		log.Println("Error viewing ApiKeys key value store")
		log.Println(err)
	}

	if name != "" { isValid = true } else { isValid = false } 

	err = db.Close()
	if err != nil {
		log.Println("Error closing key value database")
		log.Println(err)
	}

	return isValid
}

func getKVStoreRW(kvStore string) *bolt.DB {
	db, err := bolt.Open(kvStore, 0666, nil)
	if err != nil { log.Fatal(err) }
	// defer db.Close()

	return db
}

func getKVStoreRO(kvStore string) *bolt.DB {
	db, err := bolt.Open(kvStore, 0666, &bolt.Options{ReadOnly: true})
	if err != nil { log.Fatal(err) }
	// defer db.Close()

	return db
}

// Init the KV store in case we are starting empty and need some buckets made
// Called from the main program at run time...
func InitKV(kvStore string) error {

	db := getKVStoreRW(kvStore)

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("ApiKeys"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("FileBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})


	db.Close()

	return err

}
