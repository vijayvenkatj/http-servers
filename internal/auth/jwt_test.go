package auth


import (
    "testing"
    "time"

    "github.com/google/uuid"
)

const testSecret = "supersecret"

func TestGenerateAndValidateJWT(t *testing.T) {
    userID := uuid.New()

    // Create a valid token
    tokenString, err := MakeJWT(userID, testSecret, time.Minute)
    if err != nil {
        t.Fatalf("failed to generate JWT: %v", err)
    }

    // Validate token
    parsedID, err := ValidateJWT(tokenString, testSecret)
    if err != nil {
        t.Fatalf("failed to validate JWT: %v", err)
    }

    if parsedID != userID {
        t.Errorf("expected userID %v, got %v", userID, parsedID)
    }
}

func TestExpiredJWT(t *testing.T) {
    userID := uuid.New()

    // Create an already expired token
    tokenString, err := MakeJWT(userID, testSecret, -time.Minute)
    if err != nil {
        t.Fatalf("failed to generate JWT: %v", err)
    }

    // Validate should fail
    _, err = ValidateJWT(tokenString, testSecret)
    if err == nil {
        t.Errorf("expected validation error for expired token, got none")
    }
}

func TestWrongSecretJWT(t *testing.T) {
    userID := uuid.New()

    // Generate with correct secret
    tokenString, err := MakeJWT(userID, testSecret, time.Minute)
    if err != nil {
        t.Fatalf("failed to generate JWT: %v", err)
    }

    // Validate with wrong secret
    _, err = ValidateJWT(tokenString, "wrongsecret")
    if err == nil {
        t.Errorf("expected validation error for wrong secret, got none")
    }
}
