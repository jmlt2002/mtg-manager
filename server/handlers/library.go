package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"mtg-manager/server/db"

	"github.com/golang-jwt/jwt/v5"
)

// controls which cards are in a players library
func LibraryHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodPost:
		addCardToLibrary(w, r)
	case http.MethodGet:
		showLibrary(w, r)
	case http.MethodDelete:
		removeCardFromLib(w, r)
	default:
		fmt.Fprintln(w, "Invalid request method")
	}
}

func addCardToLibrary(w http.ResponseWriter, r *http.Request) {
	var card db.LibCard
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if exp, ok := (*claims)["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	username, ok := (*claims)["username"].(string)
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusUnauthorized)
		return
	}

	err = db.AddCardToLib(card, username)
	if err != nil {
		http.Error(w, "Error while trying to insert into DB", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Card added to library"})
}

func showLibrary(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if exp, ok := (*claims)["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}
	} else {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	username, ok := (*claims)["username"].(string)
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusUnauthorized)
		return
	}

	lib, err := db.GetLibrary(username)
	if err != nil {
		http.Error(w, "Somwthing went wrong while fetching from DB", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lib)
}

func removeCardFromLib(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		http.Error(w, "Invalid token format", http.StatusUnauthorized)
		return
	}

	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	username, ok := (*claims)["username"].(string)
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusUnauthorized)
		return
	}

	var request struct {
		CardID int64 `json:"card_id"`
		// here we use a pointer in case the arg is missing
		// this helps us differentiate from argument not provided and 0
		Quantity *int64 `json:"quantity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.CardID == 0 {
		http.Error(w, "Missing or invalid card ID", http.StatusBadRequest)
		return
	}

	quantity := int64(-1)
	if request.Quantity != nil {
		quantity = *request.Quantity
	}

	err = db.RemoveCardFromLib(request.CardID, quantity, username)
	if err != nil {
		http.Error(w, "Failed to remove card", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Card removed successfully"})
}
