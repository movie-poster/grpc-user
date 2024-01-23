package cryptopassword

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashAndSalt(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		return "", errors.New("error, no es una combinación de caracteres correcta")
	}
	return string(hash), nil
}

func ComparePasswords(hashedPassword string, plainPassword []byte) bool {
	byteHash := []byte(hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	return err == nil
}
