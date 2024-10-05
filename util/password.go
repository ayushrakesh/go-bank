package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("cannot hash password, %v", err)
	}
	return string(hashedPassword), nil
}

func CheckPassword(password string, hashedPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(password))
	return err
}
