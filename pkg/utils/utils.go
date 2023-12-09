package utils

import (
	"crypto/subtle"

	"golang.org/x/crypto/bcrypt"
)

func EncryPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func ComparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func SecureCompareString(source, target string) bool {
	if subtle.ConstantTimeEq(int32(len(source)), int32(len(target))) == 1 {
		return subtle.ConstantTimeCompare([]byte(source), []byte(target)) == 1
	}

	return false
}
