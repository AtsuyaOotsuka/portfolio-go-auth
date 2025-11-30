package repositories

import (
	"testing"
	"time"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/global_mock"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreateRefreshToken(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*user_refresh_tokens.*").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := NewUserRefreshTokenRepo(gdb)
	result, err := repo.CreateRefreshToken(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.RefreshToken == "" {
		t.Errorf("expected non-empty refresh token, got %q", result.RefreshToken)
	}
}

func TestCreateRefreshTokenFailDbErr(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*user_refresh_tokens.*").
		WillReturnError(sqlmock.ErrCancelled)
	mock.ExpectRollback()

	repo := NewUserRefreshTokenRepo(gdb)
	_, err := repo.CreateRefreshToken(1)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGetRefreshToken(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)

	refreshToken := &models.UserRefreshToken{
		UserID:       1,
		RefreshToken: "sample_refresh_token",
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30),
		IsUsed:       false,
	}

	rows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expires_at", "is_used"}).
		AddRow(1, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiresAt, refreshToken.IsUsed)
	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs(refreshToken.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(rows)

	defer cleanup()

	repo := NewUserRefreshTokenRepo(gdb)
	result, err := repo.getRefreshTokenl("sample_refresh_token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.RefreshToken != "sample_refresh_token" {
		t.Errorf("expected refresh token %v, got %v", "sample_refresh_token", result.RefreshToken)
	}
}

func TestGetRefreshTokenFailNotFound(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs("non_existent_token", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expires_at", "is_used"}))

	repo := NewUserRefreshTokenRepo(gdb)
	_, err := repo.getRefreshTokenl("non_existent_token")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGetRefreshTokenFailDbErr(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs("dberror_token", sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)

	repo := NewUserRefreshTokenRepo(gdb)
	_, err := repo.getRefreshTokenl("dberror_token")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGetUserByRefreshToken(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)

	user := &models.User{
		ID:           1,
		Email:        "user@example.com",
		PasswordHash: "hashed_password",
	}

	refreshToken := &models.UserRefreshToken{
		UserID:       user.ID,
		RefreshToken: "valid_refresh_token",
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30),
		IsUsed:       false,
	}

	tokenRows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expires_at", "is_used"}).
		AddRow(1, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiresAt, refreshToken.IsUsed)
	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs(refreshToken.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(tokenRows)

	userRows := sqlmock.NewRows([]string{"id", "email", "password_hash"}).
		AddRow(user.ID, user.Email, user.PasswordHash)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE id = \\?").
		WithArgs(user.ID, sqlmock.AnyArg()).
		WillReturnRows(userRows)

	defer cleanup()

	repo := NewUserRefreshTokenRepo(gdb)
	result, err := repo.GetUserByRefreshToken(refreshToken.RefreshToken)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != user.ID {
		t.Errorf("expected user ID %v, got %v", user.ID, result.ID)
	}
}

func TestGetUserByRefreshTokenUserNotFound(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)

	refreshToken := &models.UserRefreshToken{
		UserID:       1,
		RefreshToken: "valid_refresh_token",
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30),
		IsUsed:       false,
	}

	tokenRows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expires_at", "is_used"}).
		AddRow(1, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiresAt, refreshToken.IsUsed)
	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs(refreshToken.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(tokenRows)

	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE id = \\?").
		WithArgs(refreshToken.UserID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password_hash"}))

	repo := NewUserRefreshTokenRepo(gdb)
	_, err := repo.GetUserByRefreshToken(refreshToken.RefreshToken)
	if err == nil {
		t.Fatalf("expected error, got none")
	}

	defer cleanup()
}

func TestGetUserByRefreshTokenGetRefreshTokenlFail(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs("invalid_token", sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)

	repo := NewUserRefreshTokenRepo(gdb)
	_, err := repo.GetUserByRefreshToken("invalid_token")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGetUserByRefreshTokenIsUsed(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)

	refreshToken := &models.UserRefreshToken{
		UserID:       1,
		RefreshToken: "used_refresh_token",
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30),
		IsUsed:       true,
	}

	tokenRows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expires_at", "is_used"}).
		AddRow(1, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiresAt, refreshToken.IsUsed)
	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs(refreshToken.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(tokenRows)

	defer cleanup()

	repo := NewUserRefreshTokenRepo(gdb)
	_, err := repo.GetUserByRefreshToken(refreshToken.RefreshToken)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGetUserByRefreshTokenExpired(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)

	refreshToken := &models.UserRefreshToken{
		UserID:       1,
		RefreshToken: "expired_refresh_token",
		ExpiresAt:    time.Now().Add(-1 * time.Hour),
		IsUsed:       false,
	}

	tokenRows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expires_at", "is_used"}).
		AddRow(1, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiresAt, refreshToken.IsUsed)
	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs(refreshToken.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(tokenRows)

	defer cleanup()

	repo := NewUserRefreshTokenRepo(gdb)
	_, err := repo.GetUserByRefreshToken(refreshToken.RefreshToken)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestGetUserByRefreshTokenGetUserFail(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)

	refreshToken := &models.UserRefreshToken{
		UserID:       1,
		RefreshToken: "valid_refresh_token",
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30),
		IsUsed:       false,
	}

	tokenRows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expires_at", "is_used"}).
		AddRow(1, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiresAt, refreshToken.IsUsed)
	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs(refreshToken.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(tokenRows)

	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE id = \\?").
		WithArgs(refreshToken.UserID, sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)

	defer cleanup()

	repo := NewUserRefreshTokenRepo(gdb)
	_, err := repo.GetUserByRefreshToken(refreshToken.RefreshToken)
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestChangeUsed(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)

	refreshToken := &models.UserRefreshToken{
		UserID:       1,
		RefreshToken: "valid_refresh_token",
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30),
		IsUsed:       false,
	}

	tokenRows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expires_at", "is_used"}).
		AddRow(1, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiresAt, refreshToken.IsUsed)
	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs(refreshToken.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(tokenRows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user_refresh_tokens` SET").
		WithArgs(true, "192.168.0.1", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	defer cleanup()

	repo := NewUserRefreshTokenRepo(gdb)
	err := repo.ChangeUsed(refreshToken.RefreshToken, "192.168.0.1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestChangeUsedGetRefreshTokenlFail(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	defer cleanup()

	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs("invalid_token", sqlmock.AnyArg()).
		WillReturnError(sqlmock.ErrCancelled)

	repo := NewUserRefreshTokenRepo(gdb)
	err := repo.ChangeUsed("invalid_token", "192.168.0.1")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestChangeUsedFailUpdates(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	defer cleanup()

	refreshToken := &models.UserRefreshToken{
		UserID:       1,
		RefreshToken: "valid_refresh_token",
		ExpiresAt:    time.Now().Add(24 * time.Hour * 30),
		IsUsed:       false,
	}

	tokenRows := sqlmock.NewRows([]string{"id", "user_id", "refresh_token", "expires_at", "is_used"}).
		AddRow(1, refreshToken.UserID, refreshToken.RefreshToken, refreshToken.ExpiresAt, refreshToken.IsUsed)
	mock.ExpectQuery("SELECT .* FROM `user_refresh_tokens`.*WHERE refresh_token = \\?").
		WithArgs(refreshToken.RefreshToken, sqlmock.AnyArg()).
		WillReturnRows(tokenRows)

	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `user_refresh_tokens` SET").
		WithArgs(true, "192.168.0.1", sqlmock.AnyArg(), 1).
		WillReturnError(sqlmock.ErrCancelled)
	mock.ExpectRollback()

	defer cleanup()

	repo := NewUserRefreshTokenRepo(gdb)
	err := repo.ChangeUsed(refreshToken.RefreshToken, "192.168.0.1")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}
