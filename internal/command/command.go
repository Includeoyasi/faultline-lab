// Package command executes bounded commands with ownership of their process tree.
package command

import "context"

// Executor owns child processes and their bounded output for one explicit cwd/env.
type Executor struct {
	Directory string
	Env       []string
}

// Output captures both streams and the observed process exit code.
type Output struct {
	Stdout, Stderr string
	ExitCode       int
}

// Execute starts a program directly and joins its process tree before returning.
func (e Executor) Execute(ctx context.Context, program string, args []string) (Output, error) {
	inner := processExecutor{directory: e.Directory, env: e.Env}
	output, err := inner.execute(ctx, program, args)
	return Output{output.stdout, output.stderr, output.exitCode}, err
}
