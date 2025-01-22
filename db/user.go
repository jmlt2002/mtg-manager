package db

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	DateOfCreation string `json:"date_of_creation"`
}

func CreateNewUser(u User) error {
	stmt, err := Database.Prepare(`INSERT INTO users (username, password, date_of_creation) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Username, u.Password, u.DateOfCreation)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}
