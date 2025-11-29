package service

import "github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"

type CsrfSvcInterface interface {
	CreateCSRFToken(
		csrf atylabcsrf.CsrfPkgInterface,
		timestamp int64,
		secret string,
	) string
}

type CsrfSvcStruct struct{}

func NewCsrfSvcStruct() CsrfSvcInterface {
	return &CsrfSvcStruct{}
}

func (s *CsrfSvcStruct) CreateCSRFToken(
	csrf atylabcsrf.CsrfPkgInterface,
	timestamp int64,
	secret string,
) string {
	nonceStr := csrf.GenerateNonceString()
	return csrf.GenerateCSRFCookieToken(secret, timestamp, nonceStr)
}
