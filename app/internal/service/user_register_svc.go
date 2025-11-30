package service

import (
	"strings"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/public_lib/atylabencrypt"
	"gorm.io/gorm"
)

type UserRegisterSvcInterface interface {
	RegisterUser(input RegisterUserInput) (models.User, error)
}

type UserRegisterSvcStruct struct {
	db         *gorm.DB
	encryptlib atylabencrypt.EncryptPkgInterface
}

func NewUserRegisterSvc(
	db *gorm.DB,
	encryptlib atylabencrypt.EncryptPkgInterface,
) *UserRegisterSvcStruct {
	return &UserRegisterSvcStruct{
		db:         db,
		encryptlib: encryptlib,
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
		UUID:         models.UserCreateUUID(),
		Username:     input.Name,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}
