package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"gorm.io/gorm"
)

type UserRefreshTokenRepoInterface interface {
	CreateRefreshToken(userId uint) (*models.UserRefreshToken, error)
	GetUserByRefreshToken(refreshToken string) (*models.User, error)
	ChangeUsed(refreshToken string, ipAddress string) error
}

type UserRefreshTokenRepoStruct struct {
	db *gorm.DB
}

func NewUserRefreshTokenRepo(
	db *gorm.DB,
) *UserRefreshTokenRepoStruct {
	return &UserRefreshTokenRepoStruct{
		db: db,
	}
}

func (r *UserRefreshTokenRepoStruct) CreateRefreshToken(userId uint) (*models.UserRefreshToken, error) {
	model := &models.UserRefreshToken{
		UserID:       userId,
		RefreshToken: models.CreateRefreshToken(),
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30),
	}
	if err := r.db.Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (r *UserRefreshTokenRepoStruct) getRefreshTokenl(refreshToken string) (*models.UserRefreshToken, error) {
	var userRefreshToken models.UserRefreshToken
	if err := r.db.Where("refresh_token = ?", refreshToken).First(&userRefreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("refresh token not found")
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &userRefreshToken, nil
}

func (r *UserRefreshTokenRepoStruct) GetUserByRefreshToken(refreshToken string) (*models.User, error) {
	var user models.User
	userRefreshToken, err := r.getRefreshTokenl(refreshToken)
	if err != nil {
		return nil, err
	}

	if userRefreshToken.IsUsed {
		return nil, errors.New("refresh token already used")
	}

	if time.Now().After(userRefreshToken.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	if err := r.db.Where("id = ?", userRefreshToken.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user by refresh token: %w", err)
	}

	return &user, nil
}

func (r *UserRefreshTokenRepoStruct) ChangeUsed(refreshToken string, ipAddress string) error {
	userRefreshToken, err := r.getRefreshTokenl(refreshToken)
	if err != nil {
		return err
	}

	updates := map[string]any{
		"is_used": true,
		"use_ip":  ipAddress,
	}

	if err := r.db.Model(&models.UserRefreshToken{}).
		Where("id = ?", userRefreshToken.ID).
		Updates(updates).Error; err != nil {
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	return nil
}
