package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	applicationauth "github.com/TheaCOLINOT/go-project/internal/application/auth"
	applicationuser "github.com/TheaCOLINOT/go-project/internal/application/user"
	"github.com/TheaCOLINOT/go-project/internal/config"
	deliveryhttp "github.com/TheaCOLINOT/go-project/internal/delivery/http"
	"github.com/TheaCOLINOT/go-project/internal/infrastructure/persistence/postgres"
	"github.com/TheaCOLINOT/go-project/internal/infrastructure/security"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	settings, err := config.LoadDatabase()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	db, err := postgres.Connect(ctx, settings)
	cancel()
	if err != nil {
		return err
	}
	defer db.Close()

	users := postgres.NewUserRepository(db)
	tokens := postgres.NewAuthTokenRepository(db)
	authService := applicationauth.New(users, tokens, passwordHasher{})
	userService := applicationuser.New(users)

	server := &http.Server{
		Addr:              config.HTTPAddress(),
		Handler:           deliveryhttp.NewRouter(authService, userService),
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	log.Printf("HTTP server listening on %s", server.Addr)
	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server: %w", err)
		}
		return nil
	case <-ctx.Done():
	}

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()
	if err := server.Shutdown(shutdownCtx); err != nil {
		_ = server.Close()
		return fmt.Errorf("shut down HTTP server: %w", err)
	}
	if err := <-serverErrors; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server: %w", err)
	}

	return nil
}

type passwordHasher struct{}

func (passwordHasher) Hash(password string) (string, error) {
	return security.HashPassword(password)
}

func (passwordHasher) Verify(encoded, password string) bool {
	return security.VerifyPassword(encoded, password)
}
