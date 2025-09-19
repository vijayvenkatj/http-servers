package helpers

import (
	"fmt"
	"strings"
)

func ValidateChirp(chirp string) (string, error) {
	if len([]rune(chirp)) > 140 {
		return "", fmt.Errorf("chirp is too long");
	}

	temp := chirp
	words := strings.Split(temp, " ");

	for idx, word := range words {
		if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax" {
			words[idx] = "****";
		}
	}

	cleanedBody := strings.Join(words, " ");

	return cleanedBody, nil;
}