package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

// Define a struct to represent the user data
type User struct {
	Login        string `json:"login"`
	Repositories int    `json:"public_repos"`
}

var userServiceSchema graphql.Schema
var repositoryServiceURL string

func main() {
	// Define the GraphQL schema for the UserService
	var err error
	userServiceSchema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: createUserServiceQueryType(),
	})
	if err != nil {
		panic(err)
	}

	repositoryServiceURL = "http://localhost:8082/graphql" // Update with the actual URL of the RepositoryService

	// Create a new HTTP handler for the GraphQL endpoint
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), userServiceSchema)
		json.NewEncoder(w).Encode(result)
	})

	http.ListenAndServe(":8081", nil)
}

// createUserServiceQueryType creates the GraphQL query type for the UserService
func createUserServiceQueryType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type:        createUserType(),
				Description: "Get user data",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					username, _ := params.Args["username"].(string)

					// Make a GraphQL query to the RepositoryService to retrieve repository data
					repositoryResult := executeRemoteQuery(repositoryServiceURL, fmt.Sprintf(`{ repository(username: "%s", repo: "SampleProject") { repositories } }`, username))
					if repositoryResult.HasErrors() {
						return nil, repositoryResult.Errors[0]
					}

					// Combine the user data and repository data
					repositories := int(repositoryResult.Data.(map[string]interface{})["repository"].(map[string]interface{})["repositories"].(float64))
					userData := User{
						Login:        username,
						Repositories: repositories,
					}

					return userData, nil
				},
			},
		},
	})
}

// createUserType creates the GraphQL object type for the User struct
func createUserType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type: graphql.String,
			},
			"repositories": &graphql.Field{
				Type: graphql.Int,
			},
		},
	})
}

// executeQuery executes a GraphQL query against a schema
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		formattedErrors := make([]gqlerrors.FormattedError, len(result.Errors))
		for i, err := range result.Errors {
			formattedErrors[i] = gqlerrors.FormatError(err)
		}
		result.Errors = formattedErrors
	}
	return result
}

// executeRemoteQuery executes a GraphQL query against a remote schema
func executeRemoteQuery(url, query string) *graphql.Result {
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return &graphql.Result{
			Errors: gqlerrors.FormatErrors(err),
		}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.URL.RawQuery = fmt.Sprintf("query=%s", query)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &graphql.Result{
			Errors: gqlerrors.FormatErrors(err),
		}
	}
	defer resp.Body.Close()

	var result graphql.Result
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return &graphql.Result{
			Errors: gqlerrors.FormatErrors(err),
		}
	}

	return &result
}
