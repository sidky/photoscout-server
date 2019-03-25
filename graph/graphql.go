package graph

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sidky/photoscout-server/auth"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
)

type GraphQL struct {
	schema *graphql.Schema
	authenticator *auth.Authenticator
}

func NewGraphQL(schema *graphql.Schema, authenticator *auth.Authenticator) *GraphQL {
	return &GraphQL{schema: schema, authenticator: authenticator}
}

func (g *GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	token := r.Header.Get("X-Auth-Token")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	r.Context()

	uuid, err := g.authenticator.Authenticate(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if uuid == nil {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	fmt.Printf("UUID: %s\n", *uuid)
	updated := context.WithValue(r.Context(), "uuid", uuid)
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := g.schema.Exec(updated, params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
