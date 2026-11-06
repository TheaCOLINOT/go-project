package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	domainuser "github.com/TheaCOLINOT/go-project/internal/domain/user"
	"github.com/jackc/pgx/v5/pgconn"
)

// UserRepository persists users in PostgreSQL.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a PostgreSQL-backed user repository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a user and returns its database-generated fields.
func (r *UserRepository) Create(ctx context.Context, account domainuser.User) (domainuser.User, error) {
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO users (business_id, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, business_id, email, password_hash, role, is_confirmed, created_at`,
		account.BusinessID, account.Email, account.PasswordHash, account.Role,
	).Scan(
		&account.ID, &account.BusinessID, &account.Email, &account.PasswordHash,
		&account.Role, &account.IsConfirmed, &account.CreatedAt,
	)
	if err != nil {
		var postgresError *pgconn.PgError
		if errors.As(err, &postgresError) &&
			postgresError.Code == "23505" &&
			postgresError.ConstraintName == "users_email_lower_unique" {
			return domainuser.User{}, domainuser.ErrEmailAlreadyExists
		}
		return domainuser.User{}, fmt.Errorf("insert user: %w", err)
	}
	return account, nil
}

// FindByID returns the user with the given ID.
func (r *UserRepository) FindByID(ctx context.Context, id string) (domainuser.User, error) {
	account, err := scanUser(r.db.QueryRowContext(ctx, `
		SELECT id, business_id, email, password_hash, role, is_confirmed, created_at
		FROM users WHERE id = $1`, id))
	if err != nil {
		return domainuser.User{}, fmt.Errorf("find user by ID: %w", err)
	}
	return account, nil
}

// FindByEmail returns the user with the given email address, ignoring case.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domainuser.User, error) {
	account, err := scanUser(r.db.QueryRowContext(ctx, `
		SELECT id, business_id, email, password_hash, role, is_confirmed, created_at
		FROM users WHERE lower(email) = lower($1)`, email))
	if err != nil {
		return domainuser.User{}, fmt.Errorf("find user by email: %w", err)
	}
	return account, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanUser(row rowScanner) (domainuser.User, error) {
	var account domainuser.User
	err := row.Scan(
		&account.ID, &account.BusinessID, &account.Email, &account.PasswordHash,
		&account.Role, &account.IsConfirmed, &account.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domainuser.User{}, domainuser.ErrNotFound
	}
	if err != nil {
		return domainuser.User{}, err
	}
	return account, nil
}
