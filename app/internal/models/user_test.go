package models

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/public_lib/atylabencrypt"
)

func TestCreateUUID(t *testing.T) {
	uuid := UserCreateUUID()
	uuid2 := UserCreateUUID()
	uuidLength := len(uuid)
	if uuid == uuid2 {
		t.Error("Expected different UUIDs, got the same")
	}
	if uuidLength != 36 {
		t.Errorf("Expected UUID length of 36, got %d", uuidLength)
	}
	if uuid == "" {
		t.Error("Expected a valid UUID, got an empty string")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "securepassword"
	atylabEncryptPkg := atylabencrypt.NewEncryptPkg()
	hashedPassword, err := atylabEncryptPkg.CreatePasswordHash(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}
	user := User{
		PasswordHash: hashedPassword,
	}

	// Test correct password
	if err := user.VerifyPassword(password); err != nil {
		t.Errorf("Expected password to be valid, got error: %v", err)
	}

	// Test incorrect password
	if err := user.VerifyPassword("wrongpassword"); err == nil {
		t.Error("Expected password to be invalid, but got no error")
	}
}
