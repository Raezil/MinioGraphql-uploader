package main

import (
	"fmt"
	"jwt"
	"log"
	"net/http"

	prisma "db"
	"graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const defaultPort = "8080"

const (
	bucketName     = "videos"
	minioEndpoint  = "localhost:9000"
	minioAccessKey = "minioadmin"
	minioSecretKey = "minioadmin"
	useSSL         = false
)

func initializeMinioClient() (*minio.Client, error) {
	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}
	return minioClient, nil
}

func main() {
	// Initialize Prisma client
	client := prisma.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Prisma.Disconnect()
	minioClient, err := initializeMinioClient()
	if err != nil {
		log.Printf("Error initializing MinIO client: %v", err)
		return
	}
	// Initialize GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Client:      client,
		MinioClient: minioClient,
	}}))

	// Create HTTP server with JWT middleware
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", jwt.AuthMiddleware(srv))

	log.Printf("connect to http://localhost:8080/ for GraphQL playground")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
