package svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/stretchr/testify/mock"
)

type AuthSvcMock struct {
	mock.Mock
}

func (m *AuthSvcMock) Login(input service.LoginInput) (*service.AuthOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*service.AuthOutput), args.Error(1)
}

func (m *AuthSvcMock) Refresh(input service.RefreshInput) (*service.AuthOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*service.AuthOutput), args.Error(1)
}
