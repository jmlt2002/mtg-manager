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

func LoginRequest() error {
	return nil
}

func RegisterRequest() error {
	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Enter username: (leave blank to exit register)")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)

		fmt.Print("Enter password: (leave blank to exit register)")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		if username == "" || password == "" {
			return fmt.Errorf("user exit")
		}

		user := User{
			Username: username,
			Password: password,
		}

		jsonData, err := json.Marshal(user)
		if err != nil {
			return fmt.Errorf("failed to encode user data: %v", err)
		}

		url := fmt.Sprintf("%s/users", BaseURL)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			return fmt.Errorf("failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")

		fmt.Printf("Attempting to register\n")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		// TODO: make sure server sends response message
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to register, status: %d, response: %s", resp.StatusCode, string(body))
		}

		fmt.Println("User registered successfully!")

		return nil
	}
}
