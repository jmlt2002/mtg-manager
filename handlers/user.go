package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"mtg-manager/db"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodPost:
		postUser(w, r)
	case http.MethodPut:
		fmt.Fprintln(w, "put")
	case http.MethodGet:
		fmt.Fprintln(w, "get")
	case http.MethodDelete:
		fmt.Fprintln(w, "delete")
	default:
		fmt.Fprintln(w, "Invalid request method")
	}
}

func postUser(w http.ResponseWriter, r *http.Request) {
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
