package svc_mock

import (
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
	"github.com/stretchr/testify/mock"
)

type CsrfSvcMockStruct struct {
	mock.Mock
}

func (m *CsrfSvcMockStruct) CreateCSRFToken(
	csrf atylabcsrf.CsrfPkgInterface,
	timestamp int64,
	secret string,
) string {
	args := m.Called(csrf, timestamp, secret)
	return args.String(0)
}
