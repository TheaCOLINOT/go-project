package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"net/url"

	"github.com/TheaCOLINOT/go-project/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect opens and verifies a PostgreSQL connection.
func Connect(ctx context.Context, settings config.Database) (*sql.DB, error) {
	connectionURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(settings.User, settings.Password),
		Host:   net.JoinHostPort(settings.Host, fmt.Sprint(settings.Port)),
		Path:   settings.Name,
	}
	query := connectionURL.Query()
	query.Set("sslmode", settings.SSLMode)
	connectionURL.RawQuery = query.Encode()

	db, err := sql.Open("pgx", connectionURL.String())
	if err != nil {
		return nil, fmt.Errorf("open PostgreSQL connection: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping PostgreSQL: %w", err)
	}

	return db, nil
}
