package service

import (
	"fmt"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/public_lib/atylabencrypt"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/global_mock"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestRegisterUserSuccess(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*users.*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	defer cleanup()
	mock.ExpectCommit()

	input := RegisterUserInput{
		Name:     "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	encryptlibMock := new(atylabencrypt.EncryptPkgStructMock)
	encryptlibMock.On("CreatePasswordHash", input.Password).
		Return("hashedpassword123", nil)

	svc := NewUserRegisterSvc(gdb, encryptlibMock)

	user, err := svc.RegisterUser(input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Username != input.Name {
		t.Errorf("expected username %v, got %v", input.Name, user.Username)
	}
	if user.Email != input.Email {
		t.Errorf("expected email %v, got %v", input.Email, user.Email)
	}
	if user.PasswordHash != "hashedpassword123" {
		t.Errorf("expected password %v, got %v", "hashedpassword123", user.PasswordHash)
	}
	encryptlibMock.AssertExpectations(t)
}

func TestRegisterUserCreatePasswordHashError(t *testing.T) {
	gdb, _, cleanup := global_mock.NewGormWithMockError(t)
	defer cleanup()

	input := RegisterUserInput{
		Name:     "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	encryptlibMock := new(atylabencrypt.EncryptPkgStructMock)
	encryptlibMock.On("CreatePasswordHash", input.Password).
		Return("", fmt.Errorf("hash error"))

	svc := NewUserRegisterSvc(gdb, encryptlibMock)

	user, err := svc.RegisterUser(input)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if user != (models.User{}) {
		t.Errorf("expected empty user, got %v", user)
	}

	encryptlibMock.AssertExpectations(t)
}

func TestRegisterUserDBCreateError(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	defer cleanup()

	input := RegisterUserInput{
		Name:     "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	encryptlibMock := new(atylabencrypt.EncryptPkgStructMock)
	encryptlibMock.On("CreatePasswordHash", input.Password).
		Return("hashedpassword123", nil)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*users.*").
		WillReturnError(fmt.Errorf("db create error"))
	mock.ExpectRollback()

	svc := NewUserRegisterSvc(gdb, encryptlibMock)

	user, err := svc.RegisterUser(input)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if user != (models.User{}) {
		t.Errorf("expected empty user, got %v", user)
	}

	encryptlibMock.AssertExpectations(t)
}
