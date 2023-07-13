package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
)

// Define a struct to represent the repository details
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
}

// Define a struct to represent the user data
type User struct {
	Login       string       `json:"login"`
	Repositories []Repository `json:"repos"`
}

func main() {
	// Define GraphQL types for the Repository and User
	repositoryType := graphql.NewObject(graphql.ObjectConfig{
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

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"login": &graphql.Field{
				Type: graphql.String,
			},
			"repos": &graphql.Field{
				Type: graphql.NewList(repositoryType),
			},
		},
	})

	// Define the root query object
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"user": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"username": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					username := p.Args["username"].(string)

					// Make an HTTP request to the GitHub API to fetch user and repository data
					url := fmt.Sprintf("https://api.github.com/users/%s/repos", username)
					req, err := http.NewRequest("GET", url, nil)
					if err != nil {
						return nil, err
					}
					req.Header.Set("Accept", "application/vnd.github.v3+json")

					client := &http.Client{}
					resp, err := client.Do(req)
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return nil, err
					}

					var repositories []Repository
					err = json.Unmarshal(body, &repositories)
					if err != nil {
						return nil, err
					}

					user := User{
						Login:       username,
						Repositories: repositories,
					}

					return user, nil
				},
			},
		},
	})

	// Define the schema with the root query
	schemaConfig := graphql.SchemaConfig{
		Query: rootQuery,
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		fmt.Println("Failed to create schema:", err)
		return
	}

	// Define the HTTP handler for the GraphQL endpoint
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})

		if len(result.Errors) > 0 {
			errorMessages := make([]string, len(result.Errors))
			for i, err := range result.Errors {
				errorMessages[i] = err.Message
			}
			http.Error(w, strings.Join(errorMessages, "\n"), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(result)
	})

	// Start the server
	fmt.Println("Server running at http://localhost:8080/graphql")
	http.ListenAndServe(":8080", nil)
}
