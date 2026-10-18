// Command labdata prepares and verifies the lab's isolated infrastructure.
package main

import (
	"context"
	"encoding/json"
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
	"github.com/Includeoyasi/faultline-lab/internal/infrastructure"
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
		_, _ = fmt.Fprintln(errs, "usage: labdata {secrets|bootstrap|migrate|seed|checksum|smoke|reset|check-infrastructure} [flags]")
		return 2
	}
	flags := flag.NewFlagSet("labdata "+args[0], flag.ContinueOnError)
	flags.SetOutput(errs)
	if args[0] == "checksum" {
		flags.Usage = func() {
			_, _ = fmt.Fprint(errs, "Usage: labdata checksum [flags]\n\nPrint SHA-256 of migration history and Inventory product-stock data without changing it.\n\n")
			flags.PrintDefaults()
		}
	}
	config := flags.String("config", setting("FL_DATA_CONFIG", "configs/postgres/local.json"), "public database configuration")
	secrets := flags.String("secrets", setting("FL_DATA_SECRETS", ".local/postgres"), "private credentials directory")
	service := flags.String("service", "", "single migration service; empty applies all")
	project := flags.String("project", "", "isolated Compose project for explicit reset")
	confirm := flags.String("confirm", "", "repeat the project name to confirm data deletion")
	fullStack := flags.Bool("infrastructure", false, "include Kafka and Redis in secrets, smoke or reset")
	infraConfig := flags.String("infrastructure-config", setting("FL_INFRA_CONFIG", "configs/infrastructure/local.json"), "public Kafka/Redis configuration")
	infraSecrets := flags.String("infrastructure-secrets", setting("FL_INFRA_SECRETS", ".local/infrastructure"), "private Redis credentials directory")
	stage := flags.String("stage", "", "infrastructure check: before-restart or after-restart")
	checkpoint := flags.String("checkpoint", "", "acknowledged Kafka checkpoint JSON from before-restart")
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
		if *fullStack {
			err = datalab.RunInfrastructure(ctx, out)
		} else {
			err = datalab.Run(ctx, out)
		}
	case "reset":
		if *fullStack {
			err = datalab.ResetInfrastructure(ctx, *project, *confirm, out)
		} else {
			err = datalab.Reset(ctx, *project, *confirm, out)
		}
	case "check-infrastructure":
		var cfg infrastructure.Config
		cfg, err = infrastructure.Load(*infraConfig, *infraSecrets)
		if err == nil {
			operationCtx, cancel := context.WithTimeout(ctx, time.Minute)
			defer cancel()
			var result infrastructure.Checkpoint
			result, err = infrastructure.Check(operationCtx, cfg, *stage, *checkpoint)
			if err == nil {
				err = json.NewEncoder(out).Encode(result)
			}
		}
	case "secrets", "bootstrap", "migrate", "seed", "checksum":
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
			if err == nil && *fullStack {
				var infraCfg infrastructure.Config
				infraCfg, err = infrastructure.Load(*infraConfig, *infraSecrets)
				if err == nil {
					err = infrastructure.GenerateSecrets(infraCfg)
				}
			}
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
		case "checksum":
			var value string
			value, err = database.Checksum(operationCtx, cfg)
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
	if args[0] != "checksum" && args[0] != "check-infrastructure" {
		_, _ = fmt.Fprintln(out, args[0]+": PASS")
	}
	return 0
}
