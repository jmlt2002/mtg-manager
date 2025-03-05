package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type LibCard struct {
	CardID   int64  `json:"cardID"`
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}

func GetLibraryRequest(token string) error {
	url := fmt.Sprintf("%s/library", BaseURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	fmt.Println("Fetching library...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to get library: %s", string(body))
	}

	var library []LibCard
	if err := json.NewDecoder(resp.Body).Decode(&library); err != nil {
		return fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Println("Library contents:")
	for _, card := range library {
		fmt.Printf("CardID: %d, Name: %s, Quantity: %d\n", card.CardID, card.Name, card.Quantity)
	}

	return nil
}

func AddCardtoLibRequest(token string) error {
	reader := bufio.NewReader(os.Stdin)

	var card LibCard

	fmt.Print("Insert card ID: ")
	aux, _ := reader.ReadString('\n')
	aux = strings.TrimSpace(aux)

	cardID, err := strconv.ParseInt(aux, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid card ID: %w", err)
	}
	card.CardID = cardID

	fmt.Print("Insert quantity (will override value inserted before): ")
	aux, _ = reader.ReadString('\n')
	aux = strings.TrimSpace(aux)

	quantity, err := strconv.ParseInt(aux, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid quantity: %w", err)
	}
	card.Quantity = quantity

	jsonData, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("failed to encode user data: %v", err)
	}

	url := fmt.Sprintf("%s/library", BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	fmt.Println("Adding card to your library...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("operation failed: %s", string(body))
	}

	fmt.Println("Card added!")
	return nil
}

func RemoveCardfromLibRequest(token string) error {
	reader := bufio.NewReader(os.Stdin)

	var card LibCard

	fmt.Print("Insert card ID: ")
	aux, _ := reader.ReadString('\n')
	aux = strings.TrimSpace(aux)

	cardID, err := strconv.ParseInt(aux, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid card ID: %w", err)
	}
	card.CardID = cardID

	fmt.Print("Insert quantity (will subtract from current value): ")
	aux, _ = reader.ReadString('\n')
	aux = strings.TrimSpace(aux)

	if aux != "" {
		quantity, err := strconv.ParseInt(aux, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid quantity: %w", err)
		}
		card.Quantity = quantity
	}

	jsonData, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("failed to encode user data: %v", err)
	}

	url := fmt.Sprintf("%s/library", BaseURL)
	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	fmt.Println("Removing card from your library...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("operation failed: %s", string(body))
	}

	fmt.Println("Card added!")
	return nil
}
