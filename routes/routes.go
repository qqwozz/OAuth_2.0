package routes

import (
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"go-auth-app/config"
	"go-auth-app/handlers"
	"go-auth-app/middleware"
)

func RegisterRoutes(r *gin.Engine, auth *handlers.AuthHandler, pages *handlers.PageHandler, cfg *config.Config) {
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

	// Debug — only in debug mode
	if os.Getenv("GIN_MODE") != "release" {
		r.GET("/debug/env", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"AUTH0_DOMAIN":    cfg.Domain != "",
				"AUTH0_CLIENT_ID": cfg.ClientID != "",
				"AUTH0_REDIRECT":  cfg.RedirectURL,
			})
		})
	}

	// Pages
	r.GET("/", pages.Home)
	r.GET("/profile", middleware.IsAuthenticated(), pages.Profile)

	// Auth
	r.GET("/login", auth.Login)
	r.GET("/callback", auth.Callback)
	r.GET("/logout", auth.Logout)
}
