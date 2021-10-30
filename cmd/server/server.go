package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/bpeters-cmu/dns-threat-analyser/graph"
	"github.com/bpeters-cmu/dns-threat-analyser/graph/generated"
	"github.com/bpeters-cmu/dns-threat-analyser/pkg/auth"
	"github.com/bpeters-cmu/dns-threat-analyser/pkg/database"
	_ "github.com/bpeters-cmu/dns-threat-analyser/pkg/database"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	database.InitDB()

	// Using chi router as a middleware for Basic Auth
	router := chi.NewRouter()

	router.Use(auth.Basic())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	router.Handle("/graphql", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
