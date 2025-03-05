package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const BaseURL = "http://localhost:8080"

func LoginRequest() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: (leave blank to exit login): ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: (leave blank to exit login): ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return "", fmt.Errorf("user exit")
	}

	user := User{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to encode user data: %v", err)
	}

	url := fmt.Sprintf("%s/login", BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Attempting to login...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", string(body))
	}

	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	token := responseBody["token"]

	fmt.Println("Login successful!")
	return token, nil
}

func RegisterRequest() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username: (leave blank to exit register): ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("Enter password: (leave blank to exit register): ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	if username == "" || password == "" {
		return "", fmt.Errorf("user exit")
	}

	user := User{
		Username: username,
		Password: password,
	}

	jsonData, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to encode user data: %v", err)
	}

	url := fmt.Sprintf("%s/users", BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Attempting to register...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("registration failed: %s", string(body))
	}

	var responseBody map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	token, ok := responseBody["token"]
	if !ok {
		return "", fmt.Errorf("registration successful, but no token received")
	}

	fmt.Println("Registration successful!")
	return token, nil
}

func ChangePasswordRequest(token string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter old password: (leave blank to exit register): ")
	old_password, _ := reader.ReadString('\n')
	old_password = strings.TrimSpace(old_password)

	fmt.Print("Enter new password: (leave blank to exit register): ")
	new_password, _ := reader.ReadString('\n')
	new_password = strings.TrimSpace(new_password)

	if old_password == "" || new_password == "" {
		return fmt.Errorf("password change cancelled")
	}

	passwordData := map[string]string{
		"old_password": old_password,
		"new_password": new_password,
	}

	jsonData, err := json.Marshal(passwordData)
	if err != nil {
		return fmt.Errorf("failed to encode password data: %v", err)
	}

	url := fmt.Sprintf("%s/users", BaseURL)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete account: %s", string(body))
	}

	fmt.Println("Password changed successfully")
	return nil
}

func DeleteAccountRequest(token string) error {
	url := fmt.Sprintf("%s/users", BaseURL)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete account: %s", string(body))
	}

	fmt.Println("Account deleted successfully")
	return nil
}
