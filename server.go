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
	"github.com/gin-gonic/gin"
)

const defaultPort = "8080"

func main() {
	// Initialize Prisma client
	client := prisma.NewClient()
	if err := client.Prisma.Connect(); err != nil {
		log.Fatal(err)
	}
	defer client.Prisma.Disconnect()

	minioClient, err := InitMinIO()
	if err != nil {
		log.Fatalf("Error initializing MinIO: %v", err)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Client:      client,
		MinioClient: minioClient,
	}}))

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		playground.Handler("GraphQL playground", "/query").ServeHTTP(c.Writer, c.Request)
	})

	r.POST("/query", func(c *gin.Context) {
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			c.Request.ParseMultipartForm(32 << 20)
		}

		h := jwt.AuthMiddleware(srv)
		h.ServeHTTP(c.Writer, c.Request)
	})

	log.Printf("Connect to http://localhost:%s/ for GraphQL playground", defaultPort)
	log.Fatal(http.ListenAndServe(":"+defaultPort, r))
}
