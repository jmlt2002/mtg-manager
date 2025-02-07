package db

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
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

func CreateNewCard(c Card) error {
	_, err := Database.Exec(`INSERT INTO cards (name, type, cost, colors, description, isCustom) VALUES (?, ?, ?, ?, ?, ?)`,
		c.Name, c.Type, c.ManaCost, c.Colors, c.Description, c.IsCustom)
	if err != nil {
		return fmt.Errorf("failed insert card: %w", err)
	}

	return nil
}

func UpdateCard(c Card) error {
	_, err := Database.Exec(`UPDATE cards SET 
			name = ?, 
			type = ?, 
			cost = ?, 
			colors = ?, 
			description = ?, 
			isCustom = ? 
		WHERE id = ?`,
		c.Name, c.Type, c.ManaCost, c.Colors, c.Description, c.IsCustom, c.CardID)

	if err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}

	return nil
}

func GetCard(cID int64) (Card, error) {
	row := Database.QueryRow(`SELECT id, name, type, cost, colors, description, isCustom 
							   FROM cards 
							   WHERE cardid = ?`, cID)

	var card Card

	err := row.Scan(&card.CardID, &card.Name, &card.Type, &card.ManaCost, &card.Colors, &card.Description, &card.IsCustom)
	if err != nil {
		return card, fmt.Errorf("failed to get card: %w", err)
	}

	return card, nil
}
