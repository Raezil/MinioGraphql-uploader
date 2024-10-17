package minio

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinIO() (*minio.Client, error) {
	endpoint := "localhost:9000"    // e.g., "play.min.io"
	accessKeyID := "minioadmin"     // e.g., "minioadmin"
	secretAccessKey := "minioadmin" // e.g., "minioadmin"
	useSSL := false                 // Set true if using HTTPS

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	// Make sure the bucket exists, or create one
	bucketName := "uploads"
	location := "us-east-1"
	err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	return minioClient, nil
}
