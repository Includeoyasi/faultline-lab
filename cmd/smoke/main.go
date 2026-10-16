// Command smoke checks the containerized service skeleton from the host.
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

	"github.com/Includeoyasi/faultline-lab/internal/smoke"
)

func main() { os.Exit(mainCode()) }

func mainCode() int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return run(ctx, os.Args[1:], os.Stdout, os.Stderr)
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("smoke", flag.ContinueOnError)
	flags.SetOutput(stderr)
	root := flags.String("root", ".", "repository directory")
	checkConfig := flags.Bool("check-config", false, "validate the effective local Compose shutdown budgets without starting services")
	allowUnsafe := flags.Bool("allow-unsafe-shutdown", false, "explicitly allow incompatible shutdown budgets in check-config experiments")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 || (*allowUnsafe && !*checkConfig) {
		_, _ = fmt.Fprintln(stderr, "unexpected arguments; use -help")
		return 2
	}
	if *checkConfig {
		if err := smoke.CheckConfig(ctx, *root, *allowUnsafe, stdout); err != nil {
			_, _ = fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	}
	report, err := smoke.Run(ctx, smoke.Options{Root: *root, Output: stdout})
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
	}
	return report.ExitCode()
}
