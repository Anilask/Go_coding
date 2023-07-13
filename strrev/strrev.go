package main

import "fmt"
func main(){
	str :=45878
	rev :=Rev(str)
	fmt.Println(rev)
}
func Rev(str int) int {
    var rev int
	for i:=len(str)-1;i>=0;i--{
        rev = rev + string(str[i])
    }
	return rev
}
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Struct to represent a user
type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
}

var users []User

// GetUsers returns all users
func GetUsers(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(users)
}

// GetUser returns a specific user by ID
func GetUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for _, user := range users {
		if user.ID == params["id"] {
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	json.NewEncoder(w).Encode(User{})
}

// CreateUser creates a new user
func CreateUser(w http.ResponseWriter, req *http.Request) {
	var user User
	_ = json.NewDecoder(req.Body).Decode(&user)
	users = append(users, user)
	json.NewEncoder(w).Encode(users)
}

// DeleteUser deletes a user by ID
func DeleteUser(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	for index, user := range users {
		if user.ID == params["id"] {
			users = append(users[:index], users[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(users)
}

func main() {
	router := mux.NewRouter()

	// Mock data
	users = append(users, User{ID: "1", Username: "john", Email: "john@example.com"})
	users = append(users, User{ID: "2", Username: "jane", Email: "jane@example.com"})

	// API endpoints
	router.HandleFunc("/users", GetUsers).Methods("GET")
	router.HandleFunc("/users/{id}", GetUser).Methods("GET")
	router.HandleFunc("/users", CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
