package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
	_ "github.com/go-sql-driver/mysql"
)

// Define a struct to represent the repository details
type Repository struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"html_url"`
}

func main() {
	// Connect to the MySQL database
	db, err := sql.Open("mysql", "root:1234@tcp(172.23.128.1:3306)/test")
	if err != nil {
		fmt.Println("Failed to connect to the database:", err)
		return
	}
	defer db.Close()

	// Verify the database connection
	err = db.Ping()
	if err != nil {
		fmt.Println("Failed to ping the database:", err)
		return
	}

	// Create the repository table if it doesn't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS repositories (
			id INT AUTO_INCREMENT PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			url VARCHAR(255) NOT NULL
		)
	`)
	if err != nil {
		fmt.Println("Failed to create repository table:", err)
		return
	}

	// Define GraphQL types for the Repository
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

	// Define the root query object
	rootQuery := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"repository": &graphql.Field{
				Type: repositoryType,
				Args: graphql.FieldConfigArgument{
					"repositoryName": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					repositoryName := p.Args["repositoryName"].(string)

					// Make an HTTP request to the GitHub API to fetch repository data
					url := fmt.Sprintf("https://api.github.com/repos/%s", repositoryName)
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

					var repository Repository
					err = json.Unmarshal(body, &repository)
					if err != nil {
						return nil, err
					}

					// Store the repository details in the database
					stmt, err := db.Prepare("INSERT INTO repositories (name, description, url) VALUES (?, ?, ?)")
					if err != nil {
						return nil, err
					}
					defer stmt.Close()

					_, err = stmt.Exec(repository.Name, repository.Description, repository.URL)
					if err != nil {
						return nil, err
					}

					return repository, nil
				},
			},
		},
	}

	// Create a new GraphQL schema
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		fmt.Println("Failed to create schema:", err)
		return
	}

	// Define the HTTP handler for the GraphQL endpoint
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		// Parse the GraphQL query from the request body
		var requestBody struct {
			Query string `json:"query"`
		}
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		query := requestBody.Query

		// Execute the GraphQL query
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: query,
		})

		// Check for errors in the result
		if len(result.Errors) > 0 {
			errorMessages := make([]string, len(result.Errors))
			for i, err := range result.Errors {
				errorMessages[i] = err.Message
			}
			http.Error(w, strings.Join(errorMessages, "\n"), http.StatusInternalServerError)
			return
		}

		// Send the GraphQL response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// Start the server
	fmt.Println("Microservice running at http://localhost:8080/graphql")
	http.ListenAndServe(":8080", nil)
}
