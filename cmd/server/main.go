package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ncalibey/hackernews-go/internal/prisma"

	"github.com/99designs/gqlgen/handler"
	"github.com/ncalibey/hackernews-go/internal/graphql"
)

const defaultPort = "4000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	client := prisma.New(nil)
	resolver := &graphql.Resolver{Prisma: client}

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
