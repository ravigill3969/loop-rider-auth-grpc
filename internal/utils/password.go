package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	// Debug logging
	println("=== CheckPassword Debug ===")
	println("Plaintext password:", password)
	println("Password length:", len(password))
	println("Hash from DB:", hash)
	println("Hash length:", len(hash))

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		println("Error:", err.Error())
	} else {
		println("Password match: SUCCESS")
	}
	println("===========================")

	return err == nil
}
