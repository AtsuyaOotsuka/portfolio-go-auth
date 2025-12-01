package provider

import (
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/repositories"
	"github.com/AtsuyaOotsuka/portfolio-go-auth/internal/service"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabclock"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabencrypt"
	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabjwt"
)

func (p *Provider) bindAuthSvc() *service.AuthSvcStruct {
	return service.NewAuthSvc(
		repositories.NewUserRepo(p.db),
		repositories.NewUserRefreshTokenRepo(p.db),
		atylabjwt.NewJwtSvc(),
		atylabclock.NewClock(),
	)
}

func (p *Provider) bindRegisterSvc() *service.UserRegisterSvcStruct {
	return service.NewUserRegisterSvc(
		atylabencrypt.NewEncryptPkg(),
		repositories.NewUserRepo(p.db),
	)
}

func (p *Provider) bindCsrfSvc() *service.CsrfSvcStruct {
	return service.NewCsrfSvcStruct(
		atylabcsrf.NewCsrfPkgStruct(),
	)
}
