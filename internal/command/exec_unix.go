//go:build !windows

package command

import (
	"context"
	"errors"
	"os/exec"
	"syscall"
	"time"
)

func (p processExecutor) execute(ctx context.Context, program string, args []string) (commandOutput, error) {
	cmd := exec.CommandContext(ctx, program, args...)
	cmd.Dir, cmd.Env, cmd.WaitDelay = p.directory, p.env, 2*time.Second
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if errors.Is(err, syscall.ESRCH) {
			return nil
		}
		return err
	}
	stdout, stderr := &boundedBuffer{}, &boundedBuffer{}
	cmd.Stdout, cmd.Stderr = stdout, stderr
	err := cmd.Run()
	if cmd.ProcessState != nil {
		// A plugin may close its output handles and outlive a successful parent.
		// The command owns its process group even after the parent exits.
		if killErr := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); killErr != nil && !errors.Is(killErr, syscall.ESRCH) {
			err = errors.Join(err, killErr)
		}
	}
	code := -1
	if cmd.ProcessState != nil {
		code = cmd.ProcessState.ExitCode()
	}
	if stdout.truncated || stderr.truncated {
		err = errors.Join(err, errors.New("command output exceeded limit"))
	}
	return commandOutput{stdout.String(), stderr.String(), code}, errors.Join(err, ctx.Err())
}
