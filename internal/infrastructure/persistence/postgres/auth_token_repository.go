package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	domainuser "github.com/TheaCOLINOT/go-project/internal/domain/user"
)

// AuthTokenRepository persists authentication tokens in PostgreSQL.
type AuthTokenRepository struct {
	db *sql.DB
}

// NewAuthTokenRepository creates a PostgreSQL-backed token repository.
func NewAuthTokenRepository(db *sql.DB) *AuthTokenRepository {
	return &AuthTokenRepository{db: db}
}

// Create stores a token record containing only the token hash.
func (r *AuthTokenRepository) Create(ctx context.Context, token domainuser.AuthToken) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO auth_tokens (user_id, token_hash, purpose, expires_at)
		VALUES ($1, $2, $3, $4)`,
		token.UserID, token.TokenHash, token.Purpose, token.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("insert auth token: %w", err)
	}
	return nil
}

// Find retrieves a token by its hash and purpose.
func (r *AuthTokenRepository) Find(ctx context.Context, tokenHash, purpose string) (domainuser.AuthToken, error) {
	var token domainuser.AuthToken
	var usedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, `
		SELECT id, user_id, token_hash, purpose, expires_at, used_at
		FROM auth_tokens WHERE token_hash = $1 AND purpose = $2`,
		tokenHash, purpose,
	).Scan(&token.ID, &token.UserID, &token.TokenHash, &token.Purpose, &token.ExpiresAt, &usedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domainuser.AuthToken{}, domainuser.ErrNotFound
	}
	if err != nil {
		return domainuser.AuthToken{}, fmt.Errorf("find auth token: %w", err)
	}
	if usedAt.Valid {
		token.UsedAt = &usedAt.Time
	}
	return token, nil
}

// Revoke marks a still-active access token as used.
func (r *AuthTokenRepository) Revoke(ctx context.Context, tokenHash string, usedAt time.Time) error {
	result, err := r.db.ExecContext(ctx, `
		UPDATE auth_tokens SET used_at = $2
		WHERE token_hash = $1 AND purpose = 'access' AND used_at IS NULL AND expires_at > $2`,
		tokenHash, usedAt,
	)
	if err != nil {
		return fmt.Errorf("revoke auth token: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check revoked auth token: %w", err)
	}
	if rows == 0 {
		return domainuser.ErrNotFound
	}
	return nil
}
