package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	domainuser "github.com/TheaCOLINOT/go-project/internal/domain/user"
)

func TestAuthenticationInfrastructureErrorsReturnServerError(t *testing.T) {
	router := NewRouter(&failingAuthStub{}, routeUserStub{})
	request := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	request.Header.Set("Authorization", "Bearer valid-token")
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}

type failingAuthStub struct {
	routeAuthStub
}

func (*failingAuthStub) Authenticate(context.Context, string) (domainuser.User, error) {
	return domainuser.User{}, errors.New("database unavailable")
}
