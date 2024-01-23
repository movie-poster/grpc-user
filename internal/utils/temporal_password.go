package utils

import (
	cryptopassword "grpc-user/internal/crypto"

	"github.com/m1/go-generate-password/generator"
)

func CreateTemporalPassword() (string, string) {
	config := generator.Config{
		Length:                     8,
		IncludeSymbols:             false,
		IncludeNumbers:             true,
		IncludeLowercaseLetters:    true,
		IncludeUppercaseLetters:    true,
		ExcludeSimilarCharacters:   true,
		ExcludeAmbiguousCharacters: true,
	}
	g, _ := generator.New(&config)

	password, _ := g.Generate()

	hashedPassword, _ := cryptopassword.HashAndSalt([]byte(*password))

	return *password, hashedPassword

}

func ValidatePassword(pass string, hashPass string) bool {
	return cryptopassword.ComparePasswords(hashPass, []byte(pass))
}
