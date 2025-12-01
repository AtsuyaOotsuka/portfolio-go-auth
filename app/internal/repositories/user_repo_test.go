package repositories

import (
	"database/sql"
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/test_helper/mocks/global_mock"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	user := &models.User{
		PasswordHash: "hashed_password",
		Email:        "example@example.com",
	}

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*users.*").
		WillReturnResult(sqlmock.NewResult(1, 1))

	defer cleanup()
	mock.ExpectCommit()

	repo := NewUserRepo(gdb)
	if err := repo.Create(user); err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}
}

func TestCreateFailDbErr(t *testing.T) {
	user := &models.User{
		PasswordHash: "hashed_password",
		Email:        "example@example.com",
	}

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO .*users.*").
		WillReturnError(sql.ErrConnDone)

	defer cleanup()
	mock.ExpectRollback()

	repo := NewUserRepo(gdb)
	if err := repo.Create(user); err == nil {
		t.Fatalf("expected error, but got none")
	}
}

func TestUserRepoGetByEmail(t *testing.T) {
	user := &models.User{
		PasswordHash: "hashed_password",
		Email:        "example@example.com",
	}

	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	rows := sqlmock.NewRows([]string{"id", "password", "email"}).
		AddRow(1, user.PasswordHash, user.Email)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs(user.Email, sqlmock.AnyArg()).
		WillReturnRows(rows)

	defer cleanup()
	repo := NewUserRepo(gdb)
	result, err := repo.GetByEmail(user.Email)
	if err != nil {
		t.Fatalf("expected no error, but got %v", err)
	}

	if result.Email != user.Email {
		t.Errorf("expected email %v, but got %v", user.Email, result.Email)
	}
}

func TestUserRepoGetByEmailFailDbErr(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs("dberror@example.com", sqlmock.AnyArg()).
		WillReturnError(sql.ErrConnDone)
	defer cleanup()

	repo := NewUserRepo(gdb)
	_, err := repo.GetByEmail("dberror@example.com")
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
}

func TestUserRepoGetByEmailFailNotFound(t *testing.T) {
	gdb, mock, cleanup := global_mock.NewGormWithMock(t)
	mock.ExpectQuery("SELECT .* FROM `users`.*WHERE email = \\?").
		WithArgs("notfound@example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password", "email"}))
	defer cleanup()

	repo := NewUserRepo(gdb)
	_, err := repo.GetByEmail("notfound@example.com")
	if err == nil {
		t.Fatalf("expected error, but got none")
	}
}
