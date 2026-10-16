package main

import (
	"bytes"
	"context"
	"testing"
)

func TestCLIUsage(t *testing.T) {
	for _, tc := range []struct {
		args []string
		code int
	}{
		{[]string{"-help"}, 0},
		{[]string{"unexpected"}, 2},
		{[]string{"-unknown"}, 2},
		{[]string{"-allow-unsafe-shutdown"}, 2},
	} {
		var output bytes.Buffer
		if code := run(context.Background(), tc.args, &output, &output); code != tc.code {
			t.Fatalf("%v: code=%d output=%s", tc.args, code, output.String())
		}
	}
}
