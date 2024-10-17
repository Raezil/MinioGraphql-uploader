package main

import (
	"jwt"
	"log"
	"net/http"
	"strings"

	prisma "db"
	"graph"

	. "minio"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

func main() {
	// Initialize Prisma client
	client := prisma.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Prisma.Disconnect()

	// Initialize MinIO client
	minioClient, err := InitMinIO()
	if err != nil {
		log.Fatalf("Error initializing MinIO: %v", err)
	}

	// Create GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			Client:      client,
			MinioClient: minioClient,
		},
	}))

	// Handle GraphQL playground
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		playground.Handler("GraphQL playground", "/query").ServeHTTP(w, r)
	})

	// Handle GraphQL queries and file uploads with standard http
	http.HandleFunc("/query", func(w http.ResponseWriter, r *http.Request) {
		// Check if it's a multipart request for file uploads
		if strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			err := r.ParseMultipartForm(32 << 20) // Limit file size to 32MB
			if err != nil {
				http.Error(w, "Unable to parse multipart form", http.StatusBadRequest)
				return
			}
		}

		// Apply JWT middleware
		h := jwt.AuthMiddleware(srv)
		h.ServeHTTP(w, r)
	})

	// Start the server
	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, nil))
}
