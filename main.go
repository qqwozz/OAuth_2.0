package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type server struct {
	router *gin.Engine 
}

func NewServer() (*server, error) {
	router := gin.New()

	server := &server{
		router: router,
	}

	return server, nil
}

func (s *server) loginHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "calling login")
}

func (s *server) callbackHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "calling callback")
}

func main() {
	server, err := NewServer()

	if err != nil {
		log.Fatalf("could not create new server: %v", err)
	}

	server.router.GET("/ping", func(ctx *gin.Context) { 
		ctx.JSON(http.StatusOK, "pong")
	})

	server.router.GET("/login", server.loginHandler)

	server.router.GET("/callback", server.callbackHandler)

	if err := server.router.Run(); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}