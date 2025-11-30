package models

import (
	"testing"
	"time"
)

func TestCreateRefreshToken(t *testing.T) {
	token := CreateRefreshToken()
	time.Sleep(1 * time.Millisecond) // Ensure different tokens
	token2 := CreateRefreshToken()
	tokenLength := len(token)
	if token == token2 {
		t.Error("Expected different tokens, got the same")
	}
	if tokenLength != 128 {
		t.Errorf("Expected token length of 128, got %d", tokenLength)
	}
	if token == "" {
		t.Error("Expected a valid token, got an empty string")
	}
}
