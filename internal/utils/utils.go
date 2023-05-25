package utils

import (
	"crypto/md5"
	"encoding/hex"
	"net/mail"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// import (
// 	"golang.org/x/crypto/bcrypt"
// )

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPassword(password string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(password) >= 8 {
		hasMinLen = true
	}
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func HashDataSecurely(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

func CheckSum(data string) string {
	checksum := md5.Sum([]byte(data))

	return hex.EncodeToString(checksum[:])
}

func GenerateConfirmationCode(email, secret string) string {
	hash := md5.New()
	hash.Write([]byte(email))
	hash.Write([]byte(secret))

	return hex.EncodeToString(hash.Sum(nil))
}

func PackedCheckSum(data []string) string {
	checksum := md5.Sum([]byte(strings.Join(data, "")))

	return hex.EncodeToString(checksum[:])
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
