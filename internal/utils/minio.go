package utils

import (
	"fmt"
	"log"

	minio "github.com/minio/minio-go"
)

// MinioConnection Set up minio and initialize client
func MinioConnection(minioVal, portVal, accessVal, secretVal, bucketVal string, sslVal bool) *minio.Client {
	//endpoint := fmt.Sprintf("%s:%s", minioVal, portVal)
	endpoint := fmt.Sprintf("%s", minioVal)
	accessKeyID := accessVal
	secretAccessKey := secretVal
	useSSL := false
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}
	return minioClient
}
