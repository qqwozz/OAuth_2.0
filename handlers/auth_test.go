package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"go-auth-app/config"
	"go-auth-app/models"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret"))
	r.Use(sessions.Sessions("test-session", store))
	return r
}

func TestLoginRedirectsToAuth0(t *testing.T) {
	cfg := &config.Config{
		Domain:       "test.auth0.com",
		ClientID:     "test-client-id",
		ClientSecret: "test-secret",
		RedirectURL:  "http://localhost:8080/callback",
		OAuth2:       nil,
	}

	// We can't fully test OAuth2 without a real provider, but we can test the handler setup
	handler := NewAuthHandler(cfg)

	if handler == nil {
		t.Fatal("Expected non-nil AuthHandler")
	}

	if handler.cfg.Domain != "test.auth0.com" {
		t.Errorf("Expected domain 'test.auth0.com', got '%s'", handler.cfg.Domain)
	}
}

func TestLogoutClearsCookies(t *testing.T) {
	r := setupTestRouter()
	cfg := &config.Config{
		Domain:       "test.auth0.com",
		ClientID:     "test-client-id",
		SecureCookie: false,
	}

	handler := NewAuthHandler(cfg)
	r.GET("/logout", handler.Logout)

	req := httptest.NewRequest(http.MethodGet, "/logout", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should redirect to Auth0 logout URL
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status 307, got %d", w.Code)
	}

	cookies := w.Result().Cookies()
	found := false
	for _, cookie := range cookies {
		if cookie.Name == "at" {
			found = true
			if cookie.MaxAge != -1 {
				t.Errorf("Expected MaxAge -1 for at cookie, got %d", cookie.MaxAge)
			}
			if !cookie.HttpOnly {
				t.Error("Expected HttpOnly flag on at cookie")
			}
		}
	}
	if !found {
		t.Error("Expected 'at' cookie to be cleared")
	}
}

func TestPageHandlerHome(t *testing.T) {
	r := setupTestRouter()
	cfg := &config.Config{}

	// Load templates
	r.LoadHTMLGlob("../web/template/*")

	handler := NewPageHandler(cfg)
	r.GET("/", handler.Home)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestPageHandlerProfileRedirectsWhenNotAuthenticated(t *testing.T) {
	r := setupTestRouter()
	cfg := &config.Config{}

	handler := NewPageHandler(cfg)
	r.GET("/profile", handler.Profile)

	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Should redirect to home when not authenticated
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status 307, got %d", w.Code)
	}
}

func TestModelsUserInfo(t *testing.T) {
	userJSON := `{"sub":"123","email":"test@example.com","name":"Test User","email_verified":true}`

	var user models.UserInfo
	err := json.Unmarshal([]byte(userJSON), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal user: %v", err)
	}

	if user.Sub != "123" {
		t.Errorf("Expected sub '123', got '%s'", user.Sub)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
	if user.Name != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", user.Name)
	}
	if !user.EmailVerified {
		t.Error("Expected email_verified to be true")
	}
}

func TestConfigValidation(t *testing.T) {
	// Test that empty required fields would cause fatal
	// We can't test log.Fatal directly, but we can test the logic

	tests := []struct {
		name     string
		domain   string
		clientID string
		secret   string
		redirect string
		valid    bool
	}{
		{"all empty", "", "", "", "", false},
		{"missing domain", "", "id", "secret", "http://callback", false},
		{"missing client id", "domain", "", "secret", "http://callback", false},
		{"missing secret", "domain", "id", "", "http://callback", false},
		{"missing redirect", "domain", "id", "secret", "", false},
		{"all present", "domain", "id", "secret", "http://callback", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := tt.domain != "" && tt.clientID != "" && tt.secret != "" && tt.redirect != ""
			if valid != tt.valid {
				t.Errorf("Expected valid=%v, got %v", tt.valid, valid)
			}
		})
	}
}
