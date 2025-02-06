package db

import (
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateNewUser(u User) error {
	tx, err := Database.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back due to error: %v", err)
		}
	}()

	stmt, err := tx.Prepare(`INSERT INTO users (username, password, date_of_creation) VALUES (?, ?, ?)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	result, err := stmt.Exec(u.Username, u.Password, currentTime)
	if err != nil {
		return fmt.Errorf("failed to execute user insertion: %w", err)
	}

	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	createUserLibraryTable := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS lib%d (
			card_id INTEGER PRIMARY KEY NOT NULL,
			FOREIGN KEY (card_id) REFERENCES cards(id)
		);`, lastInsertId)

	if _, err := tx.Exec(createUserLibraryTable); err != nil {
		return fmt.Errorf("failed to create library for user %v: %w", lastInsertId, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func DeleteUser(u User) error {
	tx, err := Database.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back due to error: %v", err)
		}
	}()

	stmt, err := tx.Prepare(`DELETE FROM users WHERE username = ? AND password = ?`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(u.Username, u.Password)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with the given username and password")
	}

	deleteUserLibraryTable := fmt.Sprintf(`DROP TABLE IF EXISTS lib%v;`, u.Username)
	if _, err := tx.Exec(deleteUserLibraryTable); err != nil {
		return fmt.Errorf("failed to delete user's library table: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func GetUserByUsername(uname string) (User, error) {
	stmt, err := Database.Prepare(`SELECT username, password, date_of_creation FROM users WHERE username = ?`)
	if err != nil {
		return User{}, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var u User
	err = stmt.QueryRow(uname).Scan(&u.Username, &u.Password)
	if err != nil {
		return User{}, fmt.Errorf("User not found")
	}

	return u, nil
}
