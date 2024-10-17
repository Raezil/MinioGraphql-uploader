package main

import (
	"jwt"
	"log"
	"net/http"

	prisma "db"
	"graph"

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

	// Initialize GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Client: client,
	}}))

	// Create HTTP server with JWT middleware
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", jwt.AuthMiddleware(srv))

	log.Printf("connect to http://localhost:8080/ for GraphQL playground")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
