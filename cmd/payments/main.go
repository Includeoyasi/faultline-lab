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
	exitCode := service.Execute(ctx, "payments", os.Args[1:], os.LookupEnv, os.Stderr)
	stop()
	os.Exit(exitCode)
}
