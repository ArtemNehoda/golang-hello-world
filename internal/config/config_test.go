package config_test

import (
	"testing"

	"github.com/ArtemNehoda/golang-hello-world/internal/config"
)

// TestLoad_Defaults verifies that Load() returns the documented defaults when
// none of the environment variables are set.
func TestLoad_Defaults(t *testing.T) {
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("SERVER_PORT", "")

	cfg := config.Load()

	tests := []struct {
		field string
		got   string
		want  string
	}{
		{"DBHost", cfg.DBHost, "localhost"},
		{"DBPort", cfg.DBPort, "3306"},
		{"DBName", cfg.DBName, "golang_demo"},
		{"DBUser", cfg.DBUser, "golang_user"},
		{"DBPassword", cfg.DBPassword, "golang_pass"},
		{"ServerPort", cfg.ServerPort, "8080"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.field, tt.got, tt.want)
		}
	}
}

// TestLoad_FromEnv verifies that Load() picks up all values from environment
// variables when they are all set.
func TestLoad_FromEnv(t *testing.T) {
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_NAME", "mydb")
	t.Setenv("DB_USER", "admin")
	t.Setenv("DB_PASSWORD", "s3cr3t")
	t.Setenv("SERVER_PORT", "9090")

	cfg := config.Load()

	tests := []struct {
		field string
		got   string
		want  string
	}{
		{"DBHost", cfg.DBHost, "db.example.com"},
		{"DBPort", cfg.DBPort, "5432"},
		{"DBName", cfg.DBName, "mydb"},
		{"DBUser", cfg.DBUser, "admin"},
		{"DBPassword", cfg.DBPassword, "s3cr3t"},
		{"ServerPort", cfg.ServerPort, "9090"},
	}

	for _, tt := range tests {
		if tt.got != tt.want {
			t.Errorf("%s: got %q, want %q", tt.field, tt.got, tt.want)
		}
	}
}

// TestLoad_PartialEnv verifies that only the fields whose env vars are set are
// overridden; the rest still fall back to their defaults.
func TestLoad_PartialEnv(t *testing.T) {
	t.Setenv("DB_HOST", "custom-host")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("SERVER_PORT", "1234")

	cfg := config.Load()

	if cfg.DBHost != "custom-host" {
		t.Errorf("DBHost: got %q, want %q", cfg.DBHost, "custom-host")
	}
	if cfg.ServerPort != "1234" {
		t.Errorf("ServerPort: got %q, want %q", cfg.ServerPort, "1234")
	}
	// These should still be defaults.
	if cfg.DBPort != "3306" {
		t.Errorf("DBPort: got %q, want default %q", cfg.DBPort, "3306")
	}
	if cfg.DBName != "golang_demo" {
		t.Errorf("DBName: got %q, want default %q", cfg.DBName, "golang_demo")
	}
	if cfg.DBUser != "golang_user" {
		t.Errorf("DBUser: got %q, want default %q", cfg.DBUser, "golang_user")
	}
	if cfg.DBPassword != "golang_pass" {
		t.Errorf("DBPassword: got %q, want default %q", cfg.DBPassword, "golang_pass")
	}
}

// TestLoad_EmptyEnvVarFallsBackToDefault ensures that an explicitly empty
// environment variable is treated the same as unset (falls back to default).
func TestLoad_EmptyEnvVarFallsBackToDefault(t *testing.T) {
	// getEnv uses value != "" as the condition, so "" → default.
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("DB_USER", "")
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("SERVER_PORT", "")

	cfg := config.Load()

	if cfg.DBHost != "localhost" {
		t.Errorf("DBHost with empty env: got %q, want default %q", cfg.DBHost, "localhost")
	}
}

// TestLoad_ReturnsPointer confirms that Load() always returns a non-nil *Config.
func TestLoad_ReturnsPointer(t *testing.T) {
	cfg := config.Load()
	if cfg == nil {
		t.Fatal("Load() returned nil")
	}
}
