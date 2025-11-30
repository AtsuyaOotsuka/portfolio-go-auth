package models

import "testing"

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
