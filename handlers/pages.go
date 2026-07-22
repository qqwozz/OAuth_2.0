package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"go-auth-app/config"
	"go-auth-app/models"
)

var avatarClient = &http.Client{Timeout: 10 * time.Second}

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

func (h *PageHandler) Avatar(ctx *gin.Context) {
	session := sessions.Default(ctx)
	user := session.Get("user")
	if user == nil {
		ctx.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	claimsJSON, _ := json.Marshal(user)
	var u models.UserInfo
	if err := json.Unmarshal(claimsJSON, &u); err != nil || u.Picture == "" {
		ctx.Redirect(http.StatusTemporaryRedirect, "/profile")
		return
	}

	resp, err := avatarClient.Get(u.Picture)
	if err != nil {
		log.Printf("Error fetching avatar: %v", err)
		ctx.Redirect(http.StatusTemporaryRedirect, u.Picture)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ctx.Redirect(http.StatusTemporaryRedirect, u.Picture)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "image/jpeg"
	}
	ctx.Header("Content-Type", contentType)
	ctx.Header("Cache-Control", "public, max-age=3600")
	io.Copy(ctx.Writer, resp.Body)
}
