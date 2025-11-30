package atylabencrypt

import "github.com/stretchr/testify/mock"

type EncryptPkgStructMock struct {
	mock.Mock
}

func (e *EncryptPkgStructMock) CreatePasswordHash(password string) (string, error) {
	args := e.Called(password)
	return args.String(0), args.Error(1)
}
