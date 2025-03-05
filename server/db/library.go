package db

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type LibCard struct {
	CardID   int64  `json:"cardID"`
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
}

func AddCardToLib(cardID int64, quantity int64, username string) error {
	var cardName string

	err := Database.QueryRow(`SELECT name FROM cards WHERE card_id = ?`, cardID).Scan(&cardName)
	if err != nil {
		return fmt.Errorf("failed to fetch card name: %w", err)
	}

	_, err = Database.Exec(fmt.Sprintf(`INSERT INTO lib%s (card_id, card_name, quantity) VALUES (?, ?, ?)`, username),
		cardID, cardName, quantity)
	if err != nil {
		return fmt.Errorf("failed to add card to lib: %w", err)
	}

	return nil
}

func GetLibrary(username string) ([]LibCard, error) {
	rows, err := Database.Query(fmt.Sprintf("SELECT card_id, card_name, quantity FROM lib%s", username))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve library: %w", err)
	}
	defer rows.Close()

	var library []LibCard
	for rows.Next() {
		var card LibCard
		if err := rows.Scan(&card.CardID, &card.Name, &card.Quantity); err != nil {
			return nil, fmt.Errorf("failed to scan library row: %w", err)
		}
		library = append(library, card)
	}

	return library, nil
}

func RemoveCardFromLib(cID int64, q int64, username string) error {
	var err error

	if q == -1 {
		_, err = Database.Exec(fmt.Sprintf(`DELETE FROM lib%s WHERE card_id = ?`, username), cID)

	} else {
		_, err = Database.Exec(fmt.Sprintf(`
			UPDATE lib%s 
			SET quantity = quantity - ? 
			WHERE card_id = ? AND quantity > 0
		`, username), q, cID)
	}

	if err != nil {
		return fmt.Errorf("failed to remove card from lib: %w", err)
	}

	return nil
}
