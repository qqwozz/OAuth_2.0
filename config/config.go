package config

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

type Config struct {
	Domain        string
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	Port          string
	SecureCookie  bool
	SessionSecret string
	Provider      *oidc.Provider
	OAuth2        *oauth2.Config
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Domain:       os.Getenv("AUTH0_DOMAIN"),
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_REDIRECT_URL"),
	}

	if cfg.Domain == "" {
		log.Fatal("AUTH0_DOMAIN is required")
	}
	if cfg.ClientID == "" {
		log.Fatal("AUTH0_CLIENT_ID is required")
	}
	if cfg.ClientSecret == "" {
		log.Fatal("AUTH0_CLIENT_SECRET is required")
	}
	if cfg.RedirectURL == "" {
		log.Fatal("AUTH0_REDIRECT_URL is required")
	}

	// Port — из PORT env или дефолт
	cfg.Port = os.Getenv("PORT")
	if cfg.Port == "" {
		cfg.Port = ":8080"
	} else if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}

	// SecureCookie — из RedirectURL
	cfg.SecureCookie = strings.HasPrefix(cfg.RedirectURL, "https://")

	// SessionSecret — из SESSION_SECRET или генерируем
	cfg.SessionSecret = os.Getenv("SESSION_SECRET")
	if cfg.SessionSecret == "" {
		cfg.SessionSecret = generateSecret()
		log.Println("WARNING: Generated random session secret. Set SESSION_SECRET for persistence.")
	}

	// OIDC Provider — один раз
	providerURL := fmt.Sprintf("https://%s/", cfg.Domain)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	provider, err := oidc.NewProvider(ctx, providerURL)
	if err != nil {
		log.Fatalf("Failed to create OIDC provider: %v", err)
	}
	cfg.Provider = provider
	log.Printf("OIDC provider created: %s", providerURL)

	// OAuth2 config
	cfg.OAuth2 = &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
		Endpoint:     provider.Endpoint(),
	}

	return cfg
}

func generateSecret() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate secret: %v", err)
	}
	return base64.URLEncoding.EncodeToString(b)
}
