package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	_ "embed"

	"github.com/ArtemNehoda/golang-hello-world/internal/graphql/model"
	"github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

var parsedSchema *ast.Schema

//go:embed schema.graphql
var schema string

func init() {
	parsedSchema = gqlparser.MustLoadSchema(&ast.Source{Name: "schema.graphql", Input: schema})
}

type gqlRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

type gqlResponse struct {
	Data   interface{}   `json:"data,omitempty"`
	Errors gqlerror.List `json:"errors,omitempty"`
}

// NewGraphQLHandler returns an http.Handler that serves the GraphQL API.
func NewGraphQLHandler(resolver *Resolver) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(gqlResponse{
				Errors: gqlerror.List{{Message: "only POST requests are supported"}},
			})
			return
		}

		var req gqlRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(gqlResponse{
				Errors: gqlerror.List{{Message: "invalid JSON body: " + err.Error()}},
			})
			return
		}

		doc, gqlErrs := gqlparser.LoadQuery(parsedSchema, req.Query)
		if gqlErrs != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(gqlResponse{Errors: gqlErrs})
			return
		}

		data, execErrs := dispatch(r.Context(), resolver, doc)
		json.NewEncoder(w).Encode(gqlResponse{Data: data, Errors: execErrs})
	})
}

// dispatch routes the parsed GraphQL document to the appropriate resolver methods.
func dispatch(ctx context.Context, r *Resolver, doc *ast.QueryDocument) (interface{}, gqlerror.List) {
	result := map[string]interface{}{}

	for _, op := range doc.Operations {
		for _, rawSel := range op.SelectionSet {
			field, ok := rawSel.(*ast.Field)
			if !ok {
				continue
			}
			switch field.Name {
			case "messages":
				msgs, err := r.Query().Messages(ctx)
				if err != nil {
					return nil, gqlerror.List{{Message: err.Error()}}
				}

				result["messages"] = buildMessageItems(field, msgs)
			default:
				return nil, gqlerror.List{{Message: fmt.Sprintf("unknown query field: %q", field.Name)}}
			}
		}
	}

	return result, nil
}

// buildMessageItems projects the requested sub-fields from msgs into a plain map slice.
func buildMessageItems(field *ast.Field, msgs []*model.Message) []map[string]interface{} {
	items := make([]map[string]interface{}, 0, len(msgs))
	for _, m := range msgs {
		item := map[string]interface{}{}
		for _, rawSub := range field.SelectionSet {
			sub, ok := rawSub.(*ast.Field)
			if !ok {
				continue
			}
			switch sub.Name {
			case "id":
				item["id"] = m.ID
			case "content":
				item["content"] = m.Content
			case "author":
				item["author"] = m.Author
			case "createdAt":
				item["createdAt"] = m.CreatedAt
			}
		}
		items = append(items, item)
	}
	return items
}
