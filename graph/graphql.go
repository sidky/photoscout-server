package graph

import (
	"encoding/json"
	"net/http"

	graphql "github.com/graph-gophers/graphql-go"
)

type GraphQL struct {
	schema *graphql.Schema
}

func NewGraphQL(schema *graphql.Schema) *GraphQL {
	return &GraphQL{schema: schema}
}

func (g *GraphQL) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := g.schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}
