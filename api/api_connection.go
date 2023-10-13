package api

import (
	"bytes"
	"encoding/json"
	"github.com/graphql-go/graphql"
	"io"
	"net/http"
	"time"
)

type ConnectionModel struct {
	GraphQLSchema graphql.Schema
	ServerUri     string
	Token         string
}

func Connection(token string, serverUri string) (*ConnectionModel, error) {
	// Configure GraphQL schema
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"queryField": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	rootMutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"mutationField": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	schemaConfig := graphql.SchemaConfig{Query: rootQuery, Mutation: rootMutation}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		return nil, err
	}

	return &ConnectionModel{
		GraphQLSchema: schema,
		ServerUri:     serverUri,
		Token:         token,
	}, nil
}

func (c *ConnectionModel) executeRequest(query string, variables map[string]interface{}) (*[]byte, error) {
	c.ServerUri = "http://localhost:4600/coonective"

	requestBody, err := json.Marshal(map[string]interface{}{
		"query":     query,
		"variables": variables,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.ServerUri, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	// Adicione o cabeçalho "Origin" à solicitação HTTP.
	req.Header.Set("Origin", "*")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("Response body: %s\n", body, c.ServerUri)

	return &body, nil
}

func (c *ConnectionModel) Query(params string, variables map[string]interface{}) ResponseModel {
	Params(params)

	body, err := c.executeRequest(params, variables)
	if err != nil {
		apiError := ApiError([]string{err.Error()},
			"Connection", "Query", "GRAPHQL_QUERY_EXECUTE_FAILED", time.Now(), variables)
		return ResponseModel{
			Errors: []ErrorModel{apiError},
		}
	}

	return ResponseModel{
		Success: true,
		Body:    *body,
	}
}

func (c *ConnectionModel) Mutation(params string, variables map[string]interface{}) ResponseModel {

	body, err := c.executeRequest(params, variables)
	if err != nil {
		apiError := ApiError([]string{err.Error()},
			"Connection", "Mutation", "GRAPHQL_MUTATE_EXECUTE_FAILED", time.Now(), variables)
		return ResponseModel{
			Errors: []ErrorModel{apiError},
		}
	}

	return ResponseModel{
		Success: true,
		Body:    *body,
	}
}
