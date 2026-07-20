package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type server struct {
	router		 *gin.Engine 
	oauth2Config *oauth2.Config
}

func NewServer() (*server, error) {
	router := gin.New()

	oauth2Config, err := NewOauth2Config()
	if err != nil {
		return nil, fmt.Errorf("could not create new oauth2 config: %v", err)
	}

	server := &server{
		router: router,
		oauth2Config: oauth2Config,
	}

	return server, nil
}

func NewOauth2Config() (*oauth2.Config, error) {
	providerURL := fmt.Sprintf("https://" + os.Getenv("AUTH0_DOMAIN") + "/")
	provider, err := oidc.NewProvider(context.Background(), providerURL)
	if err != nil {
		return nil, fmt.Errorf("coudl not create new provider: %v", err)
	}

	oauth2Config := &oauth2.Config{
		ClientID: 		os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret:	os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL: 	os.Getenv("AUTH0_REDIRECT_URL"),
		Scopes: []string{"profile", "email", "photo"},
		Endpoint: provider.Endpoint(),
	}

	return oauth2Config, nil
}

func (s *server) loginHandler(ctx *gin.Context) {
	state, err := generateRandomString()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "could not login")
	}

	ctx.Redirect(http.StatusTemporaryRedirect, s.oauth2Config.AuthCodeURL(state))
}

func (s *server) callbackHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "calling callback")
}

func main() {
	server, err := NewServer()

	if err != nil {
		log.Fatalf("could not create new server: %v", err)
	}

	// Defining html tempate
	server.router.Static("/public", "web/static")
	server.router.LoadHTMLGlob("web/template/*")

	server.router.GET("/ping", func(ctx *gin.Context) { 
		ctx.JSON(http.StatusOK, "pong")
	})

	server.router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "home.html", nil)
	})
	
	server.router.GET("/profile", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "profile.html", nil)
	})

	server.router.GET("/login", server.loginHandler)
	server.router.GET("/callback", server.callbackHandler)

	if err := server.router.Run(); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}

func generateRandomString() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	s := base64.StdEncoding.EncodeToString(b)

	return s, nil
}