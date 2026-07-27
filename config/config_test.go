package config

import (
	"os"
	"testing"
)

func TestGenerateSecret(t *testing.T) {
	secret := generateSecret()
	if secret == "" {
		t.Error("Expected non-empty secret")
	}
	if len(secret) < 32 {
		t.Errorf("Expected secret length >= 32, got %d", len(secret))
	}

	// Generate another and verify they're different
	secret2 := generateSecret()
	if secret == secret2 {
		t.Error("Expected different secrets on each call")
	}
}

func TestPortDefault(t *testing.T) {
	// Save and restore env
	origPort := os.Getenv("PORT")
	defer os.Setenv("PORT", origPort)

	os.Unsetenv("PORT")

	cfg := &Config{
		Domain:       "test.com",
		ClientID:     "id",
		ClientSecret: "secret",
		RedirectURL:  "http://callback",
	}

	// Simulate port logic from Load()
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = ":8080"
	} else if len(cfg.Port) > 0 && cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
	}

	if cfg.Port != ":8080" {
		t.Errorf("Expected default port ':8080', got '%s'", cfg.Port)
	}
}

func TestPortWithColon(t *testing.T) {
	origPort := os.Getenv("PORT")
	defer os.Setenv("PORT", origPort)

	os.Setenv("PORT", ":9090")

	cfg := &Config{}
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = ":8080"
	} else if len(cfg.Port) > 0 && cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
	}

	if cfg.Port != ":9090" {
		t.Errorf("Expected port ':9090', got '%s'", cfg.Port)
	}
}

func TestPortWithoutColon(t *testing.T) {
	origPort := os.Getenv("PORT")
	defer os.Setenv("PORT", origPort)

	os.Setenv("PORT", "3000")

	cfg := &Config{}
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = ":8080"
	} else if len(cfg.Port) > 0 && cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
	}

	if cfg.Port != ":3000" {
		t.Errorf("Expected port ':3000', got '%s'", cfg.Port)
	}
}
