package service

import (
	"fmt"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/repo_mock"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabencrypt"
)

func TestRegisterUserSuccess(t *testing.T) {
	input := RegisterUserInput{
		Name:     "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	encryptlibMock := new(atylabencrypt.EncryptPkgStructMock)
	encryptlibMock.On("CreatePasswordHash", input.Password).
		Return("hashedpassword123", nil)

	userRepoMock := new(repo_mock.UserRepoMock)
	userRepoMock.On("Create", &models.User{
		Username:     input.Name,
		Email:        input.Email,
		PasswordHash: "hashedpassword123",
	}).Return(nil)

	svc := NewUserRegisterSvc(encryptlibMock, userRepoMock)

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
	input := RegisterUserInput{
		Name:     "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	encryptlibMock := new(atylabencrypt.EncryptPkgStructMock)
	encryptlibMock.On("CreatePasswordHash", input.Password).
		Return("", fmt.Errorf("hash error"))

	userRepoMock := new(repo_mock.UserRepoMock)

	svc := NewUserRegisterSvc(encryptlibMock, userRepoMock)

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
	input := RegisterUserInput{
		Name:     "testuser",
		Email:    "testuser@example.com",
		Password: "password123",
	}

	encryptlibMock := new(atylabencrypt.EncryptPkgStructMock)
	encryptlibMock.On("CreatePasswordHash", input.Password).
		Return("hashedpassword123", nil)

	userRepoMock := new(repo_mock.UserRepoMock)
	userRepoMock.On("Create", &models.User{
		Username:     input.Name,
		Email:        input.Email,
		PasswordHash: "hashedpassword123",
	}).Return(fmt.Errorf("db create error"))

	svc := NewUserRegisterSvc(encryptlibMock, userRepoMock)

	user, err := svc.RegisterUser(input)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	if user != (models.User{}) {
		t.Errorf("expected empty user, got %v", user)
	}

	encryptlibMock.AssertExpectations(t)
}
