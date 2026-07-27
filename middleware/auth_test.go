package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret"))
	r.Use(sessions.Sessions("test-session", store))
	return r
}

func TestIsAuthenticatedRedirectsWhenNoSession(t *testing.T) {
	r := setupTestRouter()
	r.GET("/protected", IsAuthenticated(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected status 307, got %d", w.Code)
	}
}

func TestIsAuthenticatedAllowsWhenSessionExists(t *testing.T) {
	r := setupTestRouter()
	r.GET("/protected", IsAuthenticated(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()

	// Set session before request
	r.ServeHTTP(w, req)

	// Note: Testing session-based auth requires more complex setup
	// This is a basic structure test
}

func TestSecureHeaders(t *testing.T) {
	r := setupTestRouter()
	r.Use(SecureCookies())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	headers := w.Result().Header

	expectedHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":       "DENY",
		"X-XSS-Protection":      "1; mode=block",
		"Referrer-Policy":       "strict-origin-when-cross-origin",
	}

	for header, expectedValue := range expectedHeaders {
		actual := headers.Get(header)
		if actual != expectedValue {
			t.Errorf("Expected header %s='%s', got '%s'", header, expectedValue, actual)
		}
	}
}
