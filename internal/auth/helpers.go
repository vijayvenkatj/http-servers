package auth

import (
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost);
	if err != nil {
		return "", err;
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {

	err := bcrypt.CompareHashAndPassword([]byte(hash),[]byte(password));
	if err != nil {
		return err;
	}

	return nil
}

func GetBearerToken(header http.Header) string {
    authorizationHeader := header.Get("Authorization")
    if authorizationHeader == "" {
        return ""
    }

    parts := strings.SplitN(authorizationHeader, " ", 2)
    if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
        return ""
    }

    return parts[1]
}

func GetAPIToken(header http.Header) string {
    authorizationHeader := header.Get("Authorization")
    if authorizationHeader == "" {
        return ""
    }

    parts := strings.SplitN(authorizationHeader, " ", 2)
    if len(parts) != 2 || !strings.EqualFold(parts[0], "ApiKey") {
        return ""
    }

    return parts[1]
}