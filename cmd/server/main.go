package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"

	"github.com/ncalibey/hackernews-go/internal/graphql"
	"github.com/ncalibey/hackernews-go/internal/prisma"
)

const defaultPort = "4000"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	client := prisma.New(nil)
	resolver := &graphql.Resolver{Prisma: client}
	mw := graphql.Authorization()

	http.Handle("/", mw(handler.Playground("GraphQL playground", "/query")))
	http.Handle("/query", mw(handler.GraphQL(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver}))))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
