package command

import (
	"bytes"
)

const outputLimit = 8 << 20

type commandOutput struct {
	stdout, stderr string
	exitCode       int
}

// processExecutor owns the command tree, including Docker CLI plugins.
type processExecutor struct {
	directory string
	env       []string
}

type boundedBuffer struct {
	buffer    bytes.Buffer
	truncated bool
}

func (b *boundedBuffer) Len() int       { return b.buffer.Len() }
func (b *boundedBuffer) String() string { return b.buffer.String() }

func (b *boundedBuffer) Write(p []byte) (int, error) {
	length := len(p)
	remaining := outputLimit - b.Len()
	if length > remaining {
		p = p[:remaining]
		b.truncated = true
	}
	_, _ = b.buffer.Write(p)
	return length, nil
}
