package atylabjwt

import "github.com/stretchr/testify/mock"

type JwtMock struct {
	mock.Mock
}

func (m *JwtMock) CreateJwt(config *JwtConfig) (string, error) {
	args := m.Called(config)
	return args.String(0), args.Error(1)
}
