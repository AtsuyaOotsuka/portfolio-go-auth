package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/public_lib/atylabclock"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/public_lib/atylabencrypt"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/public_lib/atylabjwt"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/funcs"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/repo_mock"
)

func TestLoginSuccess(t *testing.T) {
	funcs.WithEnv("JWT_SECRET_KEY", "testsecretkey", t, func() {
		crypt := atylabencrypt.NewEncryptPkg()

		clock := atylabclock.NewClockMock(
			time.Now(),
		)

		passwordHash, err := crypt.CreatePasswordHash("password")
		if err != nil {
			t.Fatalf("failed to hash password: %v", err)
		}

		userRepoMock := new(repo_mock.UserRepoMock)
		userRepoMock.On(
			"GetByEmail", "test@example.com",
		).Return(&models.User{
			ID:           1,
			UUID:         "test-uuid",
			Email:        "test@example.com",
			PasswordHash: passwordHash,
		}, nil)

		userRefreshTokenRepo := new(repo_mock.UserRefreshTokenRepoMock)
		userRefreshTokenRepo.On(
			"CreateRefreshToken", uint(1),
		).Return(&models.UserRefreshToken{
			ID:           1,
			UserID:       1,
			RefreshToken: "test-refresh-token",
			ExpiresAt:    clock.Now().Add(24 * time.Hour * 30),
		}, nil)

		jwtlib := &atylabjwt.JwtMock{}
		jwtlib.On(
			"CreateJwt",
			&atylabjwt.JwtConfig{
				Key:   []byte("testsecretkey"),
				Uuid:  "test-uuid",
				Email: "test@example.com",
				Exp:   clock.Now().Add(time.Hour * 1),
			},
		).Return("test-access-token", nil)

		authSvc := &AuthSvcStruct{
			userRepo:             userRepoMock,
			userRefreshTokenRepo: userRefreshTokenRepo,
			jwtlib:               jwtlib,
			clock:                clock,
		}

		input := LoginInput{
			Email:    "test@example.com",
			Password: "password",
		}

		out, err := authSvc.Login(input)
		if err != nil {
			t.Fatalf("expected no error, but got %v", err)
		}

		if out.AccessToken != "test-access-token" {
			t.Errorf("expected access token %v, but got %v", "test-access-token", out.AccessToken)
		}

		if out.RefreshToken != "test-refresh-token" {
			t.Errorf("expected refresh token %v, but got %v", "test-refresh-token", out.RefreshToken)
		}

		userRefreshTokenRepo.AssertExpectations(t)
		jwtlib.AssertExpectations(t)
		userRepoMock.AssertExpectations(t)
	})
}

func TestLoginFailInvalidPassword(t *testing.T) {
	crypt := atylabencrypt.NewEncryptPkg()

	passwordHash, err := crypt.CreatePasswordHash("password")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	userRepoMock := new(repo_mock.UserRepoMock)
	userRepoMock.On(
		"GetByEmail", "test@example.com",
	).Return(&models.User{
		ID:           1,
		UUID:         "test-uuid",
		Email:        "test@example.com",
		PasswordHash: passwordHash,
	}, nil)

	authSvc := &AuthSvcStruct{
		userRepo:             userRepoMock,
		userRefreshTokenRepo: nil,
		jwtlib:               nil,
		clock:                nil,
	}

	input := LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err = authSvc.Login(input)
	if err == nil {
		t.Fatalf("expected error, but got none")
	}

	userRepoMock.AssertExpectations(t)
}

func TestLoginFailUserNotFound(t *testing.T) {
	userRepoMock := new(repo_mock.UserRepoMock)
	userRepoMock.On(
		"GetByEmail", "test@example.com",
	).Return(&models.User{}, fmt.Errorf("user not found"))

	authSvc := &AuthSvcStruct{
		userRepo:             userRepoMock,
		userRefreshTokenRepo: nil,
		jwtlib:               nil,
		clock:                nil,
	}

	input := LoginInput{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, err := authSvc.Login(input)
	if err == nil {
		t.Fatalf("expected error, but got none")
	}

	userRepoMock.AssertExpectations(t)
}

func TestCreateResponseTokenCreateJwtFail(t *testing.T) {
	funcs.WithEnv("JWT_SECRET_KEY", "testsecretkey", t, func() {
		clock := atylabclock.NewClockMock(
			time.Now(),
		)

		user := &models.User{
			ID:           1,
			UUID:         "test-uuid",
			Email:        "test@example.com",
			PasswordHash: "hashed_password",
		}

		userRefreshTokenRepo := new(repo_mock.UserRefreshTokenRepoMock)
		jwtlib := &atylabjwt.JwtMock{}
		jwtlib.On(
			"CreateJwt",
			&atylabjwt.JwtConfig{
				Key:   []byte("testsecretkey"),
				Uuid:  "test-uuid",
				Email: "test@example.com",
				Exp:   clock.Now().Add(time.Hour * 1),
			},
		).Return("", fmt.Errorf("failed to create jwt"))

		authSvc := &AuthSvcStruct{
			userRepo:             nil,
			userRefreshTokenRepo: userRefreshTokenRepo,
			jwtlib:               jwtlib,
			clock:                clock,
		}

		_, err := authSvc.createResponseToken(user)
		if err == nil {
			t.Fatalf("expected error, but got none")
		}

		jwtlib.AssertExpectations(t)
	})
}

func TestCreateResponseTokenCreateRefreshTokenFail(t *testing.T) {
	funcs.WithEnv("JWT_SECRET_KEY", "testsecretkey", t, func() {
		clock := atylabclock.NewClockMock(
			time.Now(),
		)

		user := &models.User{
			ID:           1,
			UUID:         "test-uuid",
			Email:        "test@example.com",
			PasswordHash: "hashed_password",
		}

		userRefreshTokenRepo := new(repo_mock.UserRefreshTokenRepoMock)
		userRefreshTokenRepo.On(
			"CreateRefreshToken", uint(1),
		).Return(&models.UserRefreshToken{}, fmt.Errorf("failed to create refresh token"))

		jwtlib := &atylabjwt.JwtMock{}
		jwtlib.On(
			"CreateJwt",
			&atylabjwt.JwtConfig{
				Key:   []byte("testsecretkey"),
				Uuid:  "test-uuid",
				Email: "test@example.com",
				Exp:   clock.Now().Add(time.Hour * 1),
			},
		).Return("test-jwt", nil)

		authSvc := &AuthSvcStruct{
			userRepo:             nil,
			userRefreshTokenRepo: userRefreshTokenRepo,
			jwtlib:               jwtlib,
			clock:                clock,
		}

		_, err := authSvc.createResponseToken(user)
		if err == nil {
			t.Fatalf("expected error, but got none")
		}

		jwtlib.AssertExpectations(t)
	})
}

func TestRefresh(t *testing.T) {
	funcs.WithEnv("JWT_SECRET_KEY", "testsecretkey", t, func() {
		clock := atylabclock.NewClockMock(
			time.Now(),
		)

		user := &models.User{
			ID:    1,
			UUID:  "test-uuid",
			Email: "test@example.com",
		}
		userRefreshTokenRepo := new(repo_mock.UserRefreshTokenRepoMock)
		userRefreshTokenRepo.On(
			"GetUserByRefreshToken", "valid-refresh-token",
		).Return(user, nil)

		userRefreshTokenRepo.On(
			"CreateRefreshToken", uint(1),
		).Return(&models.UserRefreshToken{
			ID:           1,
			UserID:       1,
			RefreshToken: "new-refresh-token",
			ExpiresAt:    clock.Now().Add(24 * time.Hour * 30),
		}, nil)

		userRefreshTokenRepo.On(
			"ChangeUsed", "valid-refresh-token", "127.0.0.1",
		).Return(nil)

		jwtlib := &atylabjwt.JwtMock{}
		jwtlib.On(
			"CreateJwt",
			&atylabjwt.JwtConfig{
				Key:   []byte("testsecretkey"),
				Uuid:  "test-uuid",
				Email: "test@example.com",
				Exp:   clock.Now().Add(time.Hour * 1),
			},
		).Return("new-access-token", nil)

		authSvc := &AuthSvcStruct{
			userRepo:             nil,
			userRefreshTokenRepo: userRefreshTokenRepo,
			jwtlib:               jwtlib,
			clock:                clock,
		}

		out, err := authSvc.Refresh(RefreshInput{
			RefreshToken: "valid-refresh-token",
			IpAddress:    "127.0.0.1",
		})
		if err != nil {
			t.Fatalf("expected no error, but got: %v", err)
		}

		if out == nil {
			t.Fatal("expected output, but got nil")
		}
	})
}

func TestRefreshFailGetUserByRefreshToken(t *testing.T) {
	userRefreshTokenRepo := new(repo_mock.UserRefreshTokenRepoMock)
	userRefreshTokenRepo.On(
		"GetUserByRefreshToken", "invalid-refresh-token",
	).Return(&models.User{}, fmt.Errorf("invalid refresh token"))

	authSvc := &AuthSvcStruct{
		userRepo:             nil,
		userRefreshTokenRepo: userRefreshTokenRepo,
		jwtlib:               nil,
		clock:                nil,
	}

	_, err := authSvc.Refresh(RefreshInput{
		RefreshToken: "invalid-refresh-token",
		IpAddress:    "127.0.0.1",
	})
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
}

func TestRefreshFailChangeUsed(t *testing.T) {
	user := &models.User{
		ID:    1,
		UUID:  "test-uuid",
		Email: "test@example.com",
	}
	userRefreshTokenRepo := new(repo_mock.UserRefreshTokenRepoMock)
	userRefreshTokenRepo.On(
		"GetUserByRefreshToken", "valid-refresh-token",
	).Return(user, nil)

	userRefreshTokenRepo.On(
		"ChangeUsed", "valid-refresh-token", "127.0.0.1",
	).Return(fmt.Errorf("failed to change used status"))

	authSvc := &AuthSvcStruct{
		userRepo:             nil,
		userRefreshTokenRepo: userRefreshTokenRepo,
		jwtlib:               nil,
		clock:                nil,
	}

	_, err := authSvc.Refresh(RefreshInput{
		RefreshToken: "valid-refresh-token",
		IpAddress:    "127.0.0.1",
	})
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
}

func TestNewAuthSvc(t *testing.T) {
	userRepoMock := new(repo_mock.UserRepoMock)
	userRefreshTokenRepoMock := new(repo_mock.UserRefreshTokenRepoMock)
	jwtlibMock := &atylabjwt.JwtMock{}
	clockMock := atylabclock.NewClockMock(time.Now())

	authSvc := NewAuthSvc(
		userRepoMock,
		userRefreshTokenRepoMock,
		jwtlibMock,
		clockMock,
	)

	if authSvc.userRepo != userRepoMock {
		t.Errorf("expected userRepo to be set correctly")
	}

	if authSvc.userRefreshTokenRepo != userRefreshTokenRepoMock {
		t.Errorf("expected userRefreshTokenRepo to be set correctly")
	}

	if authSvc.jwtlib != jwtlibMock {
		t.Errorf("expected jwtlib to be set correctly")
	}

	if authSvc.clock != clockMock {
		t.Errorf("expected clock to be set correctly")
	}
}
