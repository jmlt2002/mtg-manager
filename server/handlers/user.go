package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mtg-manager/server/db"
)

// controls user's details
func UserHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodPost:
		createUser(w, r)
	case http.MethodPut:
		fmt.Fprintln(w, "put")
	case http.MethodGet:
		fmt.Fprintln(w, "get")
	case http.MethodDelete:
		deleteUser(w, r)
	default:
		fmt.Fprintln(w, "Invalid request method")
	}
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := db.CreateNewUser(user); err != nil {
		http.Error(w, "An error occurred while processing your request. Please try again later.", http.StatusInternalServerError)
		return
	}
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	var u db.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := db.DeleteUser(u); err != nil {
		http.Error(w, "An error occurred while processing your request. Please try again later.", http.StatusInternalServerError)
		return
	}
}
