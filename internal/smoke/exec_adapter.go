package smoke

import (
	"context"

	"github.com/Includeoyasi/faultline-lab/internal/command"
)

type commandOutput struct {
	stdout, stderr string
	exitCode       int
}

type executeFunc func(context.Context, string, []string) (commandOutput, error)

type processExecutor struct {
	directory string
	env       []string
}

func (p processExecutor) execute(ctx context.Context, program string, args []string) (commandOutput, error) {
	output, err := (command.Executor{Directory: p.directory, Env: p.env}).Execute(ctx, program, args)
	return commandOutput{output.Stdout, output.Stderr, output.ExitCode}, err
}
