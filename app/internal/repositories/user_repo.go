package repositories

import (
	"errors"
	"fmt"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	GetByEmail(email string) (*models.User, error)
}

type UserRepoStruct struct {
	db *gorm.DB
}

func NewUserRepo(
	db *gorm.DB,
) *UserRepoStruct {
	return &UserRepoStruct{
		db: db,
	}
}

func (r *UserRepoStruct) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}
