package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"go-auth-app/config"
	"go-auth-app/models"
)

type PageHandler struct {
	cfg *config.Config
}

func NewPageHandler(cfg *config.Config) *PageHandler {
	return &PageHandler{cfg: cfg}
}

func (h *PageHandler) Home(ctx *gin.Context) {
	session := sessions.Default(ctx)
	user := session.Get("user")
	ctx.HTML(http.StatusOK, "home.html", gin.H{
		"loggedIn": user != nil,
		"user":     user,
	})
}

func (h *PageHandler) Profile(ctx *gin.Context) {
	session := sessions.Default(ctx)
	user := session.Get("user")
	if user == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	claimsJSON, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshaling claims: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	var u models.UserInfo
	if err := json.Unmarshal(claimsJSON, &u); err != nil {
		log.Printf("Error unmarshaling claims: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	ctx.HTML(http.StatusOK, "profile.html", gin.H{
		"Profile": u,
	})
}
