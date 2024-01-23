package utils

import "strings"

func CleanString(text string) string {
	replacements := []string{
		"á", "a",
		"é", "e",
		"í", "i",
		"ó", "o",
		"ú", "u",
		"Á", "A",
		"É", "E",
		"Í", "I",
		"Ó", "O",
		"Ú", "U",
	}

	replacer := strings.NewReplacer(replacements...)
	return replacer.Replace(text)
}
