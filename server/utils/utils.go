package utils

import "golang.org/x/crypto/bcrypt"

// compares the hashed password stored in the DB with password entered by the user
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
