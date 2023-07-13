package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

// Define a struct to represent the repository details
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
}

var repositoryServiceSchema graphql.Schema

func main() {
	// Define the GraphQL schema for the RepositoryService
	var err error
	repositoryServiceSchema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: createRepositoryServiceQueryType(),
	})
	if err != nil {
		panic(err)
	}

	// Create a new HTTP handler for the GraphQL endpoint
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), repositoryServiceSchema)
		json.NewEncoder(w).Encode(result)
	})

	http.ListenAndServe(":8082", nil)
}

// createRepositoryServiceQueryType creates the GraphQL query type for the RepositoryService
func createRepositoryServiceQueryType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"repository": &graphql.Field{
				Type:        createRepositoryType(),
				Description: "Get repository data",
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"repo": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					username, _ := params.Args["username"].(string)
					repo, _ := params.Args["repo"].(string)
					return callGitHubAPI(username, repo)
				},
			},
		},
	})
}

// createRepositoryType creates the GraphQL object type for the Repository struct
func createRepositoryType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Repository",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"description": &graphql.Field{
				Type: graphql.String,
			},
			"url": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
}

// callGitHubAPI makes an HTTP request to the GitHub API to retrieve repository data
func callGitHubAPI(username, repo string) (Repository, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", username, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Repository{}, err
	}

	req.Header.Set("Authorization", "Bearer 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")

	resp, err := client.Do(req)
	if err != nil {
		return Repository{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Repository{}, fmt.Errorf("failed to retrieve repository data: %s", resp.Status)
	}

	var repository Repository
	err = json.NewDecoder(resp.Body).Decode(&repository)
	if err != nil {
		return Repository{}, err
	}

	return repository, nil
}

// executeQuery executes a GraphQL query against a schema
func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("GraphQL query execution errors: %v\n", result.Errors)
		return &graphql.Result{
			Errors: result.Errors,
		}
	}
	return result
}
