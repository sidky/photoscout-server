package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/sidky/photoscout-server/flickr"
	"github.com/sidky/photoscout-server/graph"

	graphql "github.com/graph-gophers/graphql-go"
)

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}

	return ":" + port, nil
}

func graphQL() *graph.GraphQL {
	d, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		panic(err)
	}
	flickr := flickr.FlickrApi(os.Getenv("FLICKR_API_KEY"))
	resolver := graph.NewResolver(flickr)

	schema := graphql.MustParseSchema(string(d), resolver)

	return graph.NewGraphQL(schema)
}

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("calling graphql")
	http.Handle("/graphql", graphQL())
	log.Printf("Listening on %s..\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
