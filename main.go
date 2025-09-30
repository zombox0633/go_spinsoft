package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zombox0633/go_spinsoft/src/config"
)

func main() {
	cfg := config.LoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := config.InitDatabase(ctx, cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Application
	app := config.NewApplication(cfg)

	// Graceful shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		log.Println("Gracefully shutting down...")

		// Shutdown server
		if err := app.Shutdown(); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if config.DB != nil {
			if err := config.DB.Close(shutdownCtx); err != nil {
				log.Printf("Database shutdown error: %v", err)
			} else {
				log.Println("Database connection closed")
			}
		}

		log.Println("Server stopped")
		os.Exit(0)
	}()

	log.Printf("Server starting on port %s", cfg.Port)
	if err := app.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
