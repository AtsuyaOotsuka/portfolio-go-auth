package service

import (
	"testing"

	"github.com/AtsuyaOotsuka/portfolio-go-lib/atylabcsrf"
)

func TestNewCsrfSvcStruct(t *testing.T) {
	svc := NewCsrfSvcStruct()
	_, ok := svc.(*CsrfSvcStruct)
	if !ok {
		t.Errorf("expected type *CsrfSvcStruct, got %T", svc)
	}
}

func TestCreateCSRFToken(t *testing.T) {
	csrfMock := &atylabcsrf.CsrfPkgMockStruct{}

	cvs := CsrfSvcStruct{}

	token := cvs.CreateCSRFToken(csrfMock, 1234567890, "test_secret")

	if token != "mocked_csrf_token" {
		t.Errorf("expected 'mocked_csrf_token', got '%s'", token)
	}
}
