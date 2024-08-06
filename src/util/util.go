package util

import (
	"golang.org/x/crypto/bcrypt"
)

// 比對密碼
//
//	https://gowebexamples.com/password-hashing/
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
