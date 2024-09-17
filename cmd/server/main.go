package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	config "github.com/nextlag/in_memory/configuration"
	"github.com/nextlag/in_memory/internal/server/initialize"
)

func main() {
	cfg := config.Load()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	srv, err := initialize.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv.Run(ctx)
}
