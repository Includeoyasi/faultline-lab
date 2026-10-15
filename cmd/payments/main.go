package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Includeoyasi/faultline-lab/internal/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	err := service.Run(ctx, "payments", "127.0.0.1:50053", os.LookupEnv, os.Stderr)
	stop()
	if err != nil {
		os.Exit(1)
	}
}
