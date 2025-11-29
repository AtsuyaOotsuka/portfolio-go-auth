package atylabencrypt

import "golang.org/x/crypto/bcrypt"

type EncryptPkgInterface interface {
	CreatePasswordHash(password string) (string, error)
}

type EncryptPkgStruct struct{}

func NewEncryptPkg() *EncryptPkgStruct {
	return &EncryptPkgStruct{}
}

func (e *EncryptPkgStruct) CreatePasswordHash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
