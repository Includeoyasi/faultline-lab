// Command labdata prepares and verifies the isolated PostgreSQL foundation.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/database"
	"github.com/Includeoyasi/faultline-lab/internal/datalab"
)

func main() {
	os.Exit(mainCode())
}

func mainCode() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return run(ctx, os.Args[1:], os.Stdout, os.Stderr)
}
func setting(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func run(ctx context.Context, args []string, out, errs io.Writer) int {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(errs, "usage: labdata {secrets|bootstrap|migrate|seed|fingerprint|smoke|reset} [flags]")
		return 2
	}
	flags := flag.NewFlagSet("labdata "+args[0], flag.ContinueOnError)
	flags.SetOutput(errs)
	config := flags.String("config", setting("FL_DATA_CONFIG", "configs/postgres/local.json"), "public database configuration")
	secrets := flags.String("secrets", setting("FL_DATA_SECRETS", ".local/postgres"), "private credentials directory")
	service := flags.String("service", "", "single migration service; empty applies all")
	project := flags.String("project", "", "isolated Compose project for explicit reset")
	confirm := flags.String("confirm", "", "repeat the project name to confirm data deletion")
	if err := flags.Parse(args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		return 2
	}
	var err error
	switch args[0] {
	case "smoke":
		err = datalab.Run(ctx, out)
	case "reset":
		err = datalab.Reset(ctx, *project, *confirm, out)
	case "secrets", "bootstrap", "migrate", "seed", "fingerprint":
		var cfg database.Config
		cfg, err = database.Load(*config, *secrets)
		if err != nil {
			break
		}
		operationCtx, cancel := context.WithTimeout(ctx, time.Minute)
		defer cancel()
		switch args[0] {
		case "secrets":
			err = database.GenerateSecrets(cfg)
		case "bootstrap":
			err = database.Bootstrap(operationCtx, cfg)
		case "migrate":
			selected := cfg.Services
			if *service != "" {
				var item database.Service
				item, err = cfg.Service(*service)
				selected = []database.Service{item}
			}
			if err == nil {
				for _, item := range selected {
					if err = database.Migrate(operationCtx, cfg, item); err != nil {
						break
					}
				}
			}
		case "seed":
			err = database.Seed(operationCtx, cfg)
		case "fingerprint":
			var value string
			value, err = database.Fingerprint(operationCtx, cfg)
			if err == nil {
				_, err = fmt.Fprintln(out, value)
			}
		}
	default:
		_, _ = fmt.Fprintln(errs, "unknown command; use labdata without arguments for usage")
		return 2
	}
	if err != nil {
		_, _ = fmt.Fprintln(errs, err)
		return 1
	}
	if args[0] != "fingerprint" {
		_, _ = fmt.Fprintln(out, args[0]+": PASS")
	}
	return 0
}
