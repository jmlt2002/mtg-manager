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

// controls the attributes of a certain card
func CardHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodPost:
		createCustomCard(w, r)
	case http.MethodPut:
		editCustomCard(w, r)
	case http.MethodGet:
		getCardDetails(w, r)
	case http.MethodDelete:
		fmt.Fprintln(w, "delete") // TODO:
	default:
		fmt.Fprintln(w, "Invalid request method")
	}
}

func createCustomCard(w http.ResponseWriter, r *http.Request) {
	var card db.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	card.IsCustom = true

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
	card.CreatedBy = username

	if err := db.CreateNewCard(card); err != nil {
		http.Error(w, "An error occurred while processing your request. Please try again later.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Card created successfully"})
}

func editCustomCard(w http.ResponseWriter, r *http.Request) {
	var card db.Card
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil || card.CardID == 0 {
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

	if username != card.CreatedBy {
		http.Error(w, "Can't edit card that isn't yours", http.StatusForbidden)
		return
	}

	if err := db.UpdateCard(card); err != nil {
		http.Error(w, "An error occurred while processing your request. Please try again later.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Card updated successfully"})
}

func getCardDetails(w http.ResponseWriter, r *http.Request) {
	var cID int64
	if err := json.NewDecoder(r.Body).Decode(&cID); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// no need to verify JWT, anyone can see the details of a card

	card, err := db.GetCard(cID)
	if err != nil {
		http.Error(w, "Card not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(card); err != nil {
		http.Error(w, "Failed to encode card data", http.StatusInternalServerError)
	}
}
