package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"workemailchecker/internal/api"
	"workemailchecker/internal/config"
)

func main() {
	cfg := config.Load()
	
	router := api.SetupRouter(cfg)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting WorkEmailChecker server on port %s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}