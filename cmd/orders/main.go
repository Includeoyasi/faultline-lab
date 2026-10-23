package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Includeoyasi/faultline-lab/internal/orders"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	exitCode := orders.Execute(ctx, os.Args[1:], os.LookupEnv, os.Stderr)
	stop()
	os.Exit(exitCode)
}
