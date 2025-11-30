package svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/models"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/stretchr/testify/mock"
)

type UserRegisterSvcStructMock struct {
	mock.Mock
}

func (m *UserRegisterSvcStructMock) RegisterUser(
	input service.RegisterUserInput,
) (models.User, error) {
	args := m.Called(input)
	return args.Get(0).(models.User), args.Error(1)
}
