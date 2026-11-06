package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	applicationauth "github.com/TheaCOLINOT/go-project/internal/application/auth"
	domainuser "github.com/TheaCOLINOT/go-project/internal/domain/user"
)

type authenticator interface {
	Authenticate(ctx context.Context, token string) (domainuser.User, error)
}

type userContextKey struct{}

// RequireAuth verifies a bearer token before calling the next handler.
func RequireAuth(auth authenticator, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		scheme, token, found := strings.Cut(header, " ")
		if !found || !strings.EqualFold(scheme, "Bearer") || token == "" || strings.ContainsAny(token, " \t\r\n") {
			writeUnauthorized(w)
			return
		}
		account, err := auth.Authenticate(r.Context(), token)
		if err != nil {
			if errors.Is(err, applicationauth.ErrInvalidToken) {
				writeUnauthorized(w)
			} else {
				writeInternalError(w)
			}
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey{}, account)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// UserFromContext returns the authenticated user attached by RequireAuth.
func UserFromContext(ctx context.Context) (domainuser.User, bool) {
	account, ok := ctx.Value(userContextKey{}).(domainuser.User)
	return account, ok
}

func writeUnauthorized(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{Error: "authentication required"})
}

func writeInternalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{Error: "internal server error"})
}
