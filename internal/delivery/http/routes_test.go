package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	applicationauth "github.com/TheaCOLINOT/go-project/internal/application/auth"
	domainuser "github.com/TheaCOLINOT/go-project/internal/domain/user"
)

func TestHealthRoute(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	NewRouter(nil, nil).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if got := response.Body.String(); got != "ok\n" {
		t.Fatalf("body = %q, want %q", got, "ok\n")
	}
	if got := response.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("Content-Type = %q, want %q", got, "text/plain; charset=utf-8")
	}
}

func TestAuthAndCurrentUserRoutes(t *testing.T) {
	account := domainuser.User{
		ID:           "user-id",
		BusinessID:   "usr-business-id",
		Email:        "ada@example.com",
		PasswordHash: "must-not-be-returned",
		Role:         "customer",
	}
	auth := &routeAuthStub{account: account}
	router := NewRouter(auth, routeUserStub{account: account})

	register := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(
		`{"email":"ada@example.com","password":"correct-password"}`,
	))
	register.Header.Set("Content-Type", "application/json")
	registerResponse := httptest.NewRecorder()
	router.ServeHTTP(registerResponse, register)
	if registerResponse.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", registerResponse.Code, http.StatusCreated)
	}
	if strings.Contains(registerResponse.Body.String(), "must-not-be-returned") ||
		strings.Contains(registerResponse.Body.String(), "password_hash") {
		t.Fatalf("register response exposed private fields: %s", registerResponse.Body.String())
	}

	login := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(
		`{"email":"ada@example.com","password":"correct-password"}`,
	))
	loginResponse := httptest.NewRecorder()
	router.ServeHTTP(loginResponse, login)
	if loginResponse.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginResponse.Code, http.StatusOK)
	}
	var response struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}
	if err := json.Unmarshal(loginResponse.Body.Bytes(), &response); err != nil {
		t.Fatalf("decode login response: %v", err)
	}
	if response.AccessToken != "valid-token" || response.TokenType != "Bearer" {
		t.Fatalf("login response = %#v", response)
	}

	unauthenticated := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	unauthenticatedResponse := httptest.NewRecorder()
	router.ServeHTTP(unauthenticatedResponse, unauthenticated)
	if unauthenticatedResponse.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated profile status = %d, want %d", unauthenticatedResponse.Code, http.StatusUnauthorized)
	}

	profile := httptest.NewRequest(http.MethodGet, "/api/users/me", nil)
	profile.Header.Set("Authorization", "Bearer valid-token")
	profileResponse := httptest.NewRecorder()
	router.ServeHTTP(profileResponse, profile)
	if profileResponse.Code != http.StatusOK {
		t.Fatalf("profile status = %d, want %d", profileResponse.Code, http.StatusOK)
	}
	if strings.Contains(profileResponse.Body.String(), "password_hash") {
		t.Fatalf("profile response exposed password hash: %s", profileResponse.Body.String())
	}

	logout := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	logout.Header.Set("Authorization", "Bearer valid-token")
	logoutResponse := httptest.NewRecorder()
	router.ServeHTTP(logoutResponse, logout)
	if logoutResponse.Code != http.StatusNoContent || auth.loggedOutToken != "valid-token" {
		t.Fatalf("logout = (%d, %q), want status %d and revoked token", logoutResponse.Code, auth.loggedOutToken, http.StatusNoContent)
	}
}

func TestRegisterRouteRejectsUnknownFields(t *testing.T) {
	router := NewRouter(&routeAuthStub{}, routeUserStub{})
	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(
		`{"email":"ada@example.com","password":"correct-password","role":"admin"}`,
	))
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusBadRequest)
	}
}

type routeAuthStub struct {
	account        domainuser.User
	loggedOutToken string
}

func (s *routeAuthStub) Register(_ context.Context, _, _ string) (domainuser.User, error) {
	return s.account, nil
}

func (s *routeAuthStub) Login(_ context.Context, _, _ string) (domainuser.User, string, error) {
	return s.account, "valid-token", nil
}

func (s *routeAuthStub) Authenticate(_ context.Context, token string) (domainuser.User, error) {
	if token != "valid-token" {
		return domainuser.User{}, applicationauth.ErrInvalidToken
	}
	return s.account, nil
}

func (s *routeAuthStub) Logout(_ context.Context, token string) error {
	if token != "valid-token" {
		return errors.New("unexpected token")
	}
	s.loggedOutToken = token
	return nil
}

type routeUserStub struct {
	account domainuser.User
}

func (s routeUserStub) GetByID(_ context.Context, id string) (domainuser.User, error) {
	if id != s.account.ID {
		return domainuser.User{}, domainuser.ErrNotFound
	}
	return s.account, nil
}
