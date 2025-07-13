package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := StartTaskManager(ctx)
	if err != nil {
		log.Fatalf("error while sturting up startTaskManager: %s", err.Error())
	}
}
