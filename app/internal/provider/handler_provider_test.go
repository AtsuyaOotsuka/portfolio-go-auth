package provider

import "testing"

func TestBindRegisterHandler(t *testing.T) {
	db := setupTestDB()

	provider := NewProvider(db)
	registerHandler := provider.BindRegisterHandler()

	if registerHandler == nil {
		t.Fatal("BindRegisterHandler returned nil")
	}
}

func TestBindAuthHandler(t *testing.T) {
	db := setupTestDB()

	provider := NewProvider(db)
	authHandler := provider.BindAuthHandler()

	if authHandler == nil {
		t.Fatal("BindAuthHandler returned nil")
	}
}

func TestBindCSRFHandler(t *testing.T) {
	db := setupTestDB()

	provider := NewProvider(db)
	csrfHandler := provider.BindCSRFHandler()

	if csrfHandler == nil {
		t.Fatal("BindCSRFHandler returned nil")
	}
}

func TestBindHealthCheckHandler(t *testing.T) {
	db := setupTestDB()

	provider := NewProvider(db)
	healthCheckHandler := provider.BindHealthCheckHandler()

	if healthCheckHandler == nil {
		t.Fatal("BindHealthCheckHandler returned nil")
	}
}
