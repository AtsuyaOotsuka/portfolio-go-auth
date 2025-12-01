package service

import (
	"strings"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/repositories"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabencrypt"
)

type UserRegisterSvcInterface interface {
	RegisterUser(input RegisterUserInput) (models.User, error)
}

type UserRegisterSvcStruct struct {
	encryptlib atylabencrypt.EncryptPkgInterface
	userRepo   repositories.UserRepoInterface
}

func NewUserRegisterSvc(
	encryptlib atylabencrypt.EncryptPkgInterface,
	userRepo repositories.UserRepoInterface,
) *UserRegisterSvcStruct {
	return &UserRegisterSvcStruct{
		encryptlib: encryptlib,
		userRepo:   userRepo,
	}
}

type RegisterUserInput struct {
	Name     string
	Email    string
	Password string
}

func (s *UserRegisterSvcStruct) RegisterUser(
	input RegisterUserInput,
) (models.User, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))
	hashedPassword, err := s.encryptlib.CreatePasswordHash(input.Password)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Username:     input.Name,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(&user); err != nil {
		return models.User{}, err
	}

	return user, nil
}
