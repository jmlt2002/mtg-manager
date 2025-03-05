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

type Card struct {
	CardID      int64  `json:"cardID"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	ManaCost    string `json:"cost"`
	Colors      string `json:"colors"`
	Description string `json:"description"`
	IsCustom    bool   `json:"isCustom"`
	CreatedBy   string `json:"createdBy"`
}

func CreateCustomCardRequest(token string) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Card Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Print("Card Type: ")
	cType, _ := reader.ReadString('\n')
	cType = strings.TrimSpace(cType)

	fmt.Print("Mana Cost (X2RW - for X colorless, 2 colorless, 1 read and 1 white): ")
	manaCost, _ := reader.ReadString('\n')
	manaCost = strings.TrimSpace(manaCost)

	fmt.Print("Colors: (e.g. UG - leave empty for colorless)")
	colors, _ := reader.ReadString('\n')
	colors = strings.TrimSpace(colors)

	fmt.Print("Description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	fmt.Print("Is this card custom?(y/n): ")
	isCustomInput, _ := reader.ReadString('\n')
	isCustomInput = strings.TrimSpace(isCustomInput)
	isCustom := false
	if strings.ToLower(isCustomInput) == "y" {
		isCustom = true
	}

	var card Card = Card{
		Name:        name,
		Type:        cType,
		ManaCost:    manaCost,
		Colors:      colors,
		Description: description,
		IsCustom:    isCustom,
	}

	jsonData, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("failed to encode user data: %v", err)
	}

	url := fmt.Sprintf("%s/cards", BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	fmt.Println("Creating card...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("card creation failed: %s", string(body))
	}

	fmt.Println("Card created!")
	return nil
}
