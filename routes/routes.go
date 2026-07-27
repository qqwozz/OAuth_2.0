package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"go-auth-app/config"
	"go-auth-app/handlers"
	"go-auth-app/middleware"
)

func RegisterRoutes(r *gin.Engine, auth *handlers.AuthHandler, pages *handlers.PageHandler, cfg *config.Config) {
	// Security headers
	r.Use(middleware.SecureCookies())

	// Sessions
	store := cookie.NewStore([]byte(cfg.SessionSecret))
	r.Use(sessions.Sessions("auth-sessions", store))

	// Static & templates
	r.Static("/public", "web/static")
	r.LoadHTMLGlob("web/template/*")

	// Health check
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Pages
	r.GET("/", pages.Home)
	r.GET("/profile", middleware.IsAuthenticated(), pages.Profile)
	r.GET("/avatar", middleware.IsAuthenticated(), pages.Avatar)

	// Auth
	r.GET("/login", auth.Login)
	r.GET("/callback", auth.Callback)
	r.GET("/logout", auth.Logout)
}
