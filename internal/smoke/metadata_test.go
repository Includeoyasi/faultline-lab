package smoke

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"
)

func TestMissingImageMetadataPreservesFunctionalResult(t *testing.T) {
	f, options := fixture(t)
	r, err := run(t.Context(), options, func(ctx context.Context, program string, args []string) (commandOutput, error) {
		if args[0] == "image" && args[1] == "inspect" {
			return commandOutput{stderr: "metadata-secret"}, errors.New("metadata-secret")
		}
		return f.execute(ctx, program, args)
	})
	if err != nil || r.Result != "PASS" || r.Cleanup != "PASS" || r.Environment.Services["orders"].Image.Status != "UNKNOWN" {
		t.Fatalf("metadata changed functional outcome: %s %v", r.Result, err)
	}
	if strings.Contains(r.Error, "metadata-secret") {
		t.Fatal("raw metadata error leaked")
	}
	verifySavedReport(t, options, r)
}

func TestCancellationDuringMetadataCleansAndStopsStages(t *testing.T) {
	f, options := fixture(t)
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	r, err := run(ctx, options, func(child context.Context, program string, args []string) (commandOutput, error) {
		if args[0] == "info" && strings.Contains(args[len(args)-1], `"cpus"`) {
			if _, ok := child.Deadline(); !ok {
				t.Fatal("metadata query lacks deadline")
			}
			cancel()
			return commandOutput{}, child.Err()
		}
		return f.execute(child, program, args)
	})
	if err == nil || !r.Cancelled || r.Result != "FAIL" || r.Cleanup != "PASS" || f.owned || f.images || !f.foreign || slices.Contains(f.operations, "exec") {
		t.Fatalf("cancelled metadata run did not stop and clean: %s %v", r.Result, err)
	}
	verifySavedReport(t, options, r)
}
