package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"go-auth-app/config"
	"go-auth-app/utils"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	state, err := utils.GenerateRandomString()
	if err != nil {
		log.Printf("Error generating state: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate state"})
		return
	}

	session := sessions.Default(ctx)
	session.Set("state", state)
	if err := session.Save(); err != nil {
		log.Printf("Error saving session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not save session"})
		return
	}

	ctx.Redirect(http.StatusTemporaryRedirect, h.cfg.OAuth2.AuthCodeURL(state))
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	ctx.SetCookie("at", "", -1, "/", "", false, h.cfg.SecureCookie)
	ctx.SetCookie("auth-sessions", "", -1, "/", "", false, h.cfg.SecureCookie)

	logoutURL := fmt.Sprintf("https://%s/v2/logout", h.cfg.Domain)

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	returnTo := scheme + "://" + ctx.Request.Host

	params := url.Values{}
	params.Add("returnTo", returnTo)
	params.Add("client_id", h.cfg.ClientID)

	ctx.Redirect(http.StatusTemporaryRedirect, logoutURL+"?"+params.Encode())
}

func (h *AuthHandler) Callback(ctx *gin.Context) {
	session := sessions.Default(ctx)

	// Validate state
	state := session.Get("state")
	queryState := ctx.Query("state")
	if state != queryState {
		log.Printf("State mismatch: session=%v, query=%v", state, queryState)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid state param"})
		return
	}
	session.Delete("state")
	_ = session.Save()

	// Get code
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing code parameter"})
		return
	}

	// Exchange code for token
	token, err := h.cfg.OAuth2.Exchange(ctx, code)
	if err != nil {
		log.Printf("Error exchanging code: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not exchange code"})
		return
	}
	if !token.Valid() {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
		return
	}

	// Save access token
	ctx.SetCookie("at", token.AccessToken, 3600, "/", "", false, h.cfg.SecureCookie)

	// Verify ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "no id_token in response"})
		return
	}

	verifier := h.cfg.Provider.Verifier(&oidc.Config{ClientID: h.cfg.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Printf("Error verifying id token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not verify id token"})
		return
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		log.Printf("Error parsing claims: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not parse claims"})
		return
	}

	session.Set("user", claims)
	if err := session.Save(); err != nil {
		log.Printf("Error saving session: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not save session"})
		return
	}

	log.Printf("User authenticated: %v", claims["email"])
	ctx.Redirect(http.StatusTemporaryRedirect, "/profile")
}
