package config

import (
	"testing"
)

func TestHTTPAddress(t *testing.T) {
	t.Setenv("HTTP_ADDR", "")
	if got := HTTPAddress(); got != ":8080" {
		t.Fatalf("HTTPAddress() = %q, want %q", got, ":8080")
	}

	t.Setenv("HTTP_ADDR", "127.0.0.1:9090")
	if got := HTTPAddress(); got != "127.0.0.1:9090" {
		t.Fatalf("HTTPAddress() = %q, want %q", got, "127.0.0.1:9090")
	}
}

func TestLoadDatabaseDefaults(t *testing.T) {
	for _, key := range []string{
		"POSTGRES_HOST",
		"POSTGRES_PORT",
		"POSTGRES_DB",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_SSLMODE",
	} {
		t.Setenv(key, "")
	}

	database, err := LoadDatabase()
	if err != nil {
		t.Fatalf("LoadDatabase() error = %v", err)
	}

	if database.Host != "localhost" || database.Port != 5432 ||
		database.Name != "ecommerce" || database.User != "ecommerce" ||
		database.Password != "dev_password_change_me" || database.SSLMode != "disable" {
		t.Fatalf("LoadDatabase() = %#v, want development defaults", database)
	}
}

func TestLoadDatabaseUsesEnvironment(t *testing.T) {
	t.Setenv("POSTGRES_HOST", "db.example")
	t.Setenv("POSTGRES_PORT", "5433")
	t.Setenv("POSTGRES_DB", "store")
	t.Setenv("POSTGRES_USER", "store_user")
	t.Setenv("POSTGRES_PASSWORD", "secret")
	t.Setenv("POSTGRES_SSLMODE", "require")

	database, err := LoadDatabase()
	if err != nil {
		t.Fatalf("LoadDatabase() error = %v", err)
	}

	if database.Host != "db.example" || database.Port != 5433 ||
		database.Name != "store" || database.User != "store_user" ||
		database.Password != "secret" || database.SSLMode != "require" {
		t.Fatalf("LoadDatabase() = %#v, want environment settings", database)
	}
}

func TestLoadDatabaseRejectsInvalidPort(t *testing.T) {
	for _, port := range []string{"not-a-port", "0", "65536"} {
		t.Run(port, func(t *testing.T) {
			t.Setenv("POSTGRES_PORT", port)
			if _, err := LoadDatabase(); err == nil {
				t.Fatal("LoadDatabase() error = nil, want invalid port error")
			}
		})
	}
}
