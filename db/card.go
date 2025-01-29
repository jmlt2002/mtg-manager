package db

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Card struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	ManaCost    string `json:"cost"`
	Colors      string `json:"colors"`
	Description string `json:"description"`
	IsCustom    int64  `json:"isCustom"`
}

func CreateNewCard(c Card) error {
	_, err := Database.Exec(`INSERT INTO cards (name, type, cost, colors, description, isCustom) VALUES (?, ?, ?, ?, ?, ?)`,
		c.Name, c.Type, c.ManaCost, c.Colors, c.Description, c.IsCustom)
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return nil
}
