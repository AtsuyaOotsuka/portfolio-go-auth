package repo_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/stretchr/testify/mock"
)

type UserRefreshTokenRepoMock struct {
	mock.Mock
}

func (m *UserRefreshTokenRepoMock) CreateRefreshToken(userId uint) (*models.UserRefreshToken, error) {
	args := m.Called(userId)
	return args.Get(0).(*models.UserRefreshToken), args.Error(1)
}

func (m *UserRefreshTokenRepoMock) GetUserByRefreshToken(refreshToken string) (*models.User, error) {
	args := m.Called(refreshToken)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *UserRefreshTokenRepoMock) ChangeUsed(refreshToken string, ipAddress string) error {
	args := m.Called(refreshToken, ipAddress)
	return args.Error(0)
}
