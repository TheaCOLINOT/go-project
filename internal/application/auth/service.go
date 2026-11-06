package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	domainuser "github.com/TheaCOLINOT/go-project/internal/domain/user"
)

const (
	accessTokenPurpose  = "access"
	accessTokenLifetime = 24 * time.Hour
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired access token")
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrInvalidPassword    = errors.New("password must be between 8 and 256 bytes")
)

// TokenRepository stores hashed access tokens.
type TokenRepository interface {
	Create(ctx context.Context, token domainuser.AuthToken) error
	Find(ctx context.Context, tokenHash, purpose string) (domainuser.AuthToken, error)
	Revoke(ctx context.Context, tokenHash string, usedAt time.Time) error
}

// PasswordHasher hashes and verifies account passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(encoded, password string) bool
}

// Service implements account registration and access-token authentication.
type Service struct {
	users     domainuser.Repository
	tokens    TokenRepository
	passwords PasswordHasher
	now       func() time.Time
}

// New creates an authentication service.
func New(users domainuser.Repository, tokens TokenRepository, passwords PasswordHasher) *Service {
	return &Service{users: users, tokens: tokens, passwords: passwords, now: time.Now}
}

// Register creates a customer account. Authentication is obtained separately through Login.
func (s *Service) Register(ctx context.Context, email, password string) (domainuser.User, error) {
	email, err := normalizeEmail(email)
	if err != nil {
		return domainuser.User{}, err
	}
	if len(password) < 8 || len(password) > 256 {
		return domainuser.User{}, ErrInvalidPassword
	}
	passwordHash, err := s.passwords.Hash(password)
	if err != nil {
		return domainuser.User{}, fmt.Errorf("hash password: %w", err)
	}
	businessID, err := newBusinessID()
	if err != nil {
		return domainuser.User{}, fmt.Errorf("generate business ID: %w", err)
	}
	created, err := s.users.Create(ctx, domainuser.User{
		BusinessID:   businessID,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         "customer",
	})
	if err != nil {
		return domainuser.User{}, fmt.Errorf("register user: %w", err)
	}
	return created, nil
}

// Login verifies credentials and returns a one-day opaque access token.
func (s *Service) Login(ctx context.Context, email, password string) (domainuser.User, string, error) {
	email, err := normalizeEmail(email)
	if err != nil || len(password) > 256 {
		return domainuser.User{}, "", ErrInvalidCredentials
	}
	account, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, domainuser.ErrNotFound) {
		return domainuser.User{}, "", ErrInvalidCredentials
	}
	if err != nil {
		return domainuser.User{}, "", fmt.Errorf("find user for login: %w", err)
	}
	if !s.passwords.Verify(account.PasswordHash, password) {
		return domainuser.User{}, "", ErrInvalidCredentials
	}

	rawToken, tokenHash, err := newAccessToken()
	if err != nil {
		return domainuser.User{}, "", fmt.Errorf("generate access token: %w", err)
	}
	now := s.now()
	if err := s.tokens.Create(ctx, domainuser.AuthToken{
		UserID:    account.ID,
		TokenHash: tokenHash,
		Purpose:   accessTokenPurpose,
		ExpiresAt: now.Add(accessTokenLifetime),
	}); err != nil {
		return domainuser.User{}, "", fmt.Errorf("store access token: %w", err)
	}
	return account, rawToken, nil
}

// Authenticate resolves an unexpired, unrevoked access token to its user.
func (s *Service) Authenticate(ctx context.Context, rawToken string) (domainuser.User, error) {
	if !validAccessToken(rawToken) {
		return domainuser.User{}, ErrInvalidToken
	}
	tokenHash := hashToken(rawToken)
	token, err := s.tokens.Find(ctx, tokenHash, accessTokenPurpose)
	if errors.Is(err, domainuser.ErrNotFound) {
		return domainuser.User{}, ErrInvalidToken
	}
	if err != nil {
		return domainuser.User{}, fmt.Errorf("find access token: %w", err)
	}
	if token.UsedAt != nil || !s.now().Before(token.ExpiresAt) {
		return domainuser.User{}, ErrInvalidToken
	}
	account, err := s.users.FindByID(ctx, token.UserID)
	if err != nil {
		if errors.Is(err, domainuser.ErrNotFound) {
			return domainuser.User{}, ErrInvalidToken
		}
		return domainuser.User{}, fmt.Errorf("find authenticated user: %w", err)
	}
	return account, nil
}

// Logout revokes the supplied access token.
func (s *Service) Logout(ctx context.Context, rawToken string) error {
	if !validAccessToken(rawToken) {
		return ErrInvalidToken
	}
	if err := s.tokens.Revoke(ctx, hashToken(rawToken), s.now()); err != nil {
		if errors.Is(err, domainuser.ErrNotFound) {
			return ErrInvalidToken
		}
		return fmt.Errorf("revoke access token: %w", err)
	}
	return nil
}

func normalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	if len(email) > 320 {
		return "", ErrInvalidEmail
	}
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email || !strings.Contains(email, "@") {
		return "", ErrInvalidEmail
	}
	return strings.ToLower(email), nil
}

func newBusinessID() (string, error) {
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "", err
	}
	id[6] = id[6]&0x0f | 0x40
	id[8] = id[8]&0x3f | 0x80
	encoded := hex.EncodeToString(id[:])
	return "usr_" + encoded[:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:], nil
}

func newAccessToken() (string, string, error) {
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return "", "", err
	}
	raw := base64.RawURLEncoding.EncodeToString(random)
	return raw, hashToken(raw), nil
}

func validAccessToken(token string) bool {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	return err == nil && len(decoded) == 32
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
