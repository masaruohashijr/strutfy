package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/masaru/strutfy-api/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	application, err := app.New(ctx)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	if err := application.Run(ctx); err != nil {
		log.Fatalf("server stopped with error: %v", err)
	}
}
