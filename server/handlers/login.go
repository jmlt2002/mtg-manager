package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"mtg-manager/server/db"
	"mtg-manager/server/utils"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	switch method {
	case http.MethodPost:
		loginUser(w, r)
	default:
		fmt.Fprintln(w, "Invalid request method")
	}
}

var secretKey = "SUPER_SECRET_KEY"

func generateJWToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(secretKey)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var input db.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	user, err := db.GetUserByUsername(input.Username)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	valid := utils.CheckPassword(user.Password, input.Password)
	if !valid {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := generateJWToken(input.Username)
	if err != nil {
		http.Error(w, "Could not generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
