package main

import (
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"go-auth-app/config"
	"go-auth-app/handlers"
	"go-auth-app/routes"
)

func main() {
	gob.Register(map[string]interface{}{})

	cfg := config.Load()

	r := gin.Default()

	auth := handlers.NewAuthHandler(cfg)
	pages := handlers.NewPageHandler(cfg)

	routes.RegisterRoutes(r, auth, pages, cfg)

	srv := &http.Server{
		Addr:    cfg.Port,
		Handler: r,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down...")
		srv.Shutdown(context.Background())
	}()

	log.Printf("Server starting on %s", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed: %v", err)
	}
	log.Println("Server stopped")
}
