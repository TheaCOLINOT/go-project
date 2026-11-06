package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	applicationauth "github.com/TheaCOLINOT/go-project/internal/application/auth"
	"github.com/TheaCOLINOT/go-project/internal/delivery/http/middleware"
	domainuser "github.com/TheaCOLINOT/go-project/internal/domain/user"
)

type authUseCases interface {
	Register(ctx context.Context, email, password string) (domainuser.User, error)
	Login(ctx context.Context, email, password string) (domainuser.User, string, error)
	Authenticate(ctx context.Context, token string) (domainuser.User, error)
	Logout(ctx context.Context, token string) error
}

type userUseCases interface {
	GetByID(ctx context.Context, id string) (domainuser.User, error)
}

// NewRouter creates the HTTP routes for the application.
func NewRouter(authService authUseCases, userService userUseCases) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("POST /api/auth/register", func(w http.ResponseWriter, r *http.Request) {
		var request credentialsRequest
		if err := decodeJSON(w, r, &request); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		account, err := authService.Register(r.Context(), request.Email, request.Password)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, toUserResponse(account))
	})

	mux.HandleFunc("POST /api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		var request credentialsRequest
		if err := decodeJSON(w, r, &request); err != nil {
			writeError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		account, token, err := authService.Login(r.Context(), request.Email, request.Password)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, loginResponse{
			User:        toUserResponse(account),
			AccessToken: token,
			TokenType:   "Bearer",
		})
	})

	mux.Handle("GET /api/users/me", middleware.RequireAuth(authService, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		account, ok := middleware.UserFromContext(r.Context())
		if !ok {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		account, err := userService.GetByID(r.Context(), account.ID)
		if err != nil {
			writeServiceError(w, err)
			return
		}
		writeJSON(w, http.StatusOK, toUserResponse(account))
	})))

	mux.Handle("POST /api/auth/logout", middleware.RequireAuth(authService, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := bearerToken(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		if err := authService.Logout(r.Context(), token); err != nil {
			writeServiceError(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})))

	return mux
}

type credentialsRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	User        publicUserResponse `json:"user"`
	AccessToken string             `json:"access_token"`
	TokenType   string             `json:"token_type"`
}

type publicUserResponse struct {
	ID          string `json:"id"`
	BusinessID  string `json:"business_id"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	IsConfirmed bool   `json:"is_confirmed"`
}

func toUserResponse(account domainuser.User) publicUserResponse {
	return publicUserResponse{
		ID:          account.ID,
		BusinessID:  account.BusinessID,
		Email:       account.Email,
		Role:        account.Role,
		IsConfirmed: account.IsConfirmed,
	}
}

func decodeJSON(w http.ResponseWriter, r *http.Request, destination any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return fmt.Errorf("request body must contain one JSON value")
	}
	return nil
}

func bearerToken(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	scheme, token, found := strings.Cut(header, " ")
	if !found || !strings.EqualFold(scheme, "Bearer") || token == "" || strings.ContainsAny(token, " \t\r\n") {
		return "", fmt.Errorf("missing bearer token")
	}
	return token, nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, applicationauth.ErrInvalidCredentials),
		errors.Is(err, applicationauth.ErrInvalidToken):
		writeError(w, http.StatusUnauthorized, "invalid credentials")
	case errors.Is(err, applicationauth.ErrInvalidEmail),
		errors.Is(err, applicationauth.ErrInvalidPassword):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, domainuser.ErrEmailAlreadyExists):
		writeError(w, http.StatusConflict, "email already exists")
	case errors.Is(err, domainuser.ErrNotFound):
		writeError(w, http.StatusNotFound, "user not found")
	default:
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, struct {
		Error string `json:"error"`
	}{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
