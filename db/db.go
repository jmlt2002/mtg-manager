package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func InitDB(filepath string) error {
	var err error
	Database, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err := createInitialTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	return nil
}

func createInitialTables() error {
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
		date_of_creation TEXT NOT NULL
	);`

	createCardsTable := `
	CREATE TABLE IF NOT EXISTS cards (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		type TEXT,
		mana_cost TEXT,
		rarity TEXT,
		colors TEXT,
		description TEXT,
		user_id INTEGER,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);`

	if _, err := Database.Exec(createUsersTable); err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}
	if _, err := Database.Exec(createCardsTable); err != nil {
		return fmt.Errorf("failed to create cards table: %v", err)
	}

	log.Println("Database tables created successfully.")
	return nil
}

// (to be used only for testing purposes, theoretically this is a server so should be running all the time)
func DestroyDB() error {
	err := Database.Close()
	if err != nil {
		return fmt.Errorf("failed to destroy db")
	}

	log.Println("Database destroyed succesfully.")
	return nil
}
