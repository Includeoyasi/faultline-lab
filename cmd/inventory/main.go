package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Includeoyasi/faultline-lab/internal/inventory"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	exitCode := inventory.Execute(ctx, os.Args[1:], os.LookupEnv, os.Stderr)
	stop()
	os.Exit(exitCode)
}
