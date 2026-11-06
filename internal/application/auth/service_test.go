package auth

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	domainuser "github.com/TheaCOLINOT/go-project/internal/domain/user"
)

func TestRegisterLoginAuthenticateAndLogout(t *testing.T) {
	users := &memoryUsers{byEmail: make(map[string]domainuser.User)}
	tokens := &memoryTokens{byHash: make(map[string]domainuser.AuthToken)}
	service := New(users, tokens, testHasher{})
	now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return now }

	account, err := service.Register(context.Background(), "  Ada@Example.com ", "correct-password")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if account.Email != "ada@example.com" || account.Role != "customer" || account.PasswordHash == "correct-password" {
		t.Fatalf("Register() account = %#v", account)
	}

	authenticated, rawToken, err := service.Login(context.Background(), "ADA@example.com", "correct-password")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if authenticated.ID != account.ID {
		t.Fatalf("Login() user ID = %q, want %q", authenticated.ID, account.ID)
	}
	tokenHash := hashToken(rawToken)
	if strings.Contains(tokens.byHash[tokenHash].TokenHash, rawToken) || tokenHash == rawToken {
		t.Fatal("stored token was not hashed")
	}

	resolved, err := service.Authenticate(context.Background(), rawToken)
	if err != nil || resolved.ID != account.ID {
		t.Fatalf("Authenticate() = (%#v, %v), want user %q", resolved, err, account.ID)
	}
	if err := service.Logout(context.Background(), rawToken); err != nil {
		t.Fatalf("Logout() error = %v", err)
	}
	if _, err := service.Authenticate(context.Background(), rawToken); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("Authenticate() after logout error = %v, want ErrInvalidToken", err)
	}
}

func TestRegisterAndLoginRejectInvalidInput(t *testing.T) {
	service := New(&memoryUsers{byEmail: make(map[string]domainuser.User)}, &memoryTokens{byHash: make(map[string]domainuser.AuthToken)}, testHasher{})

	if _, err := service.Register(context.Background(), "not-an-email", "correct-password"); !errors.Is(err, ErrInvalidEmail) {
		t.Fatalf("Register() error = %v, want ErrInvalidEmail", err)
	}
	if _, err := service.Register(context.Background(), "ada@example.com", "short"); !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("Register() error = %v, want ErrInvalidPassword", err)
	}
	if _, _, err := service.Login(context.Background(), "ada@example.com", "wrong-password"); !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want ErrInvalidCredentials", err)
	}
}

type testHasher struct{}

func (testHasher) Hash(password string) (string, error) { return "test:" + password, nil }
func (testHasher) Verify(encoded, password string) bool { return encoded == "test:"+password }

type memoryUsers struct {
	byEmail map[string]domainuser.User
	nextID  int
}

func (r *memoryUsers) Create(_ context.Context, account domainuser.User) (domainuser.User, error) {
	if _, exists := r.byEmail[account.Email]; exists {
		return domainuser.User{}, domainuser.ErrEmailAlreadyExists
	}
	r.nextID++
	account.ID = string(rune('0' + r.nextID))
	r.byEmail[account.Email] = account
	return account, nil
}

func (r *memoryUsers) FindByID(_ context.Context, id string) (domainuser.User, error) {
	for _, account := range r.byEmail {
		if account.ID == id {
			return account, nil
		}
	}
	return domainuser.User{}, domainuser.ErrNotFound
}

func (r *memoryUsers) FindByEmail(_ context.Context, email string) (domainuser.User, error) {
	account, exists := r.byEmail[email]
	if !exists {
		return domainuser.User{}, domainuser.ErrNotFound
	}
	return account, nil
}

type memoryTokens struct {
	byHash map[string]domainuser.AuthToken
}

func (r *memoryTokens) Create(_ context.Context, token domainuser.AuthToken) error {
	r.byHash[token.TokenHash] = token
	return nil
}

func (r *memoryTokens) Find(_ context.Context, tokenHash, _ string) (domainuser.AuthToken, error) {
	token, exists := r.byHash[tokenHash]
	if !exists {
		return domainuser.AuthToken{}, domainuser.ErrNotFound
	}
	return token, nil
}

func (r *memoryTokens) Revoke(_ context.Context, tokenHash string, usedAt time.Time) error {
	token, exists := r.byHash[tokenHash]
	if !exists || token.UsedAt != nil {
		return domainuser.ErrNotFound
	}
	token.UsedAt = &usedAt
	r.byHash[tokenHash] = token
	return nil
}
