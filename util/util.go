package util

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

//Striper is a function for trimming whitespaces
func Striper(str string) *string {
	str = strings.TrimSpace(str)
	if str == "" {
		return nil
	}
	return &str
}

//HashPassword function hashes password with bcrypt algorithm as Cost value and return hashed string value with an error
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

//CheckPasswordHash function checks two inputs and returns TRUE if matches
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
