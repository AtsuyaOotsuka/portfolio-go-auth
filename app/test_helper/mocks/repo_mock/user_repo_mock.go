package repo_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (r *UserRepoMock) GetByEmail(email string) (*models.User, error) {
	args := r.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (r *UserRepoMock) Create(user *models.User) error {
	args := r.Called(user)
	return args.Error(0)
}
