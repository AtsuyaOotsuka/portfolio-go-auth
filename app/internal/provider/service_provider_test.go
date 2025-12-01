package provider

import (
	"testing"

	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	return &gorm.DB{}
}

func TestBindAuthSvc(t *testing.T) {
	db := setupTestDB()

	provider := NewProvider(db)
	authSvc := provider.bindAuthSvc()

	if authSvc == nil {
		t.Fatal("BindAuthSvc returned nil")
	}
}

func TestBindRegisterSvc(t *testing.T) {
	db := setupTestDB()

	provider := NewProvider(db)
	registerSvc := provider.bindRegisterSvc()

	if registerSvc == nil {
		t.Fatal("BindRegisterSvc returned nil")
	}
}

func TestBindCsrfSvc(t *testing.T) {
	db := setupTestDB()

	provider := NewProvider(db)
	csrfSvc := provider.bindCsrfSvc()

	if csrfSvc == nil {
		t.Fatal("BindCsrfSvc returned nil")
	}
}
