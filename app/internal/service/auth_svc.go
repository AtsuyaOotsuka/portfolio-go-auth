package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/repositories"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/public_lib/atylabjwt"
)

type AuthSvcInterface interface {
	Login(input LoginInput) (*AuthOutput, error)
	Refresh(input RefreshInput) (*AuthOutput, error)
}

type AuthSvcStruct struct {
	userRepo             repositories.UserRepoInterface
	userRefreshTokenRepo repositories.UserRefreshTokenRepoInterface
	jwtlib               atylabjwt.JwtSvcInterface
}

func NewAuthSvc(
	userRepo repositories.UserRepoInterface,
	userRefreshTokenRepo repositories.UserRefreshTokenRepoInterface,
	jwtlib atylabjwt.JwtSvcInterface,
) *AuthSvcStruct {
	return &AuthSvcStruct{
		userRepo:             userRepo,
		userRefreshTokenRepo: userRefreshTokenRepo,
		jwtlib:               jwtlib,
	}
}

type AuthOutput struct {
	AccessToken  string
	RefreshToken string
}

type LoginInput struct {
	Email    string
	Password string
}

func (s *AuthSvcStruct) Login(input LoginInput) (*AuthOutput, error) {
	email := strings.TrimSpace(strings.ToLower(input.Email))

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("invalid email: %w", err) // 本来は曖昧にするが、学習目的のため、分ける
	}

	// パスワード検証
	if err := user.VerifyPassword(input.Password); err != nil {
		return nil, fmt.Errorf("invalid password: %w", err) // 本来は曖昧にするが、学習目的のため、分ける
	}

	return s.createResponseToken(user)
}

func (s *AuthSvcStruct) createResponseToken(user *models.User) (*AuthOutput, error) {
	// jwtを発行
	jwt, err := s.jwtlib.CreateJwt(&atylabjwt.JwtConfig{
		Key:   []byte(os.Getenv("JWT_SECRET_KEY")),
		Uuid:  user.UUID,
		Email: user.Email,
		Exp:   time.Now().Add(time.Hour * 1),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create jwt: %w", err)
	}

	refreshToken, err := s.userRefreshTokenRepo.CreateRefreshToken(user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &AuthOutput{
		AccessToken:  jwt,
		RefreshToken: refreshToken.RefreshToken,
	}, nil
}

type RefreshInput struct {
	RefreshToken string
	IpAddress    string
}

func (s *AuthSvcStruct) Refresh(input RefreshInput) (*AuthOutput, error) {
	refreshTokenRecord, err := s.userRefreshTokenRepo.GetUserByRefreshToken(input.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	if err := s.userRefreshTokenRepo.ChangeUsed(input.RefreshToken, input.IpAddress); err != nil {
		return nil, fmt.Errorf("failed to change used refresh token: %w", err)
	}

	return s.createResponseToken(refreshTokenRecord)
}
