package runmeta

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func syntheticHome(directory string) string {
	return strings.Join([]string{"", directory, "private"}, "/")
}

func dockerFixture(t *testing.T, backend, actual, requested, failure string) Docker {
	t.Helper()
	id := strings.Repeat("a", 64)
	image := "sha256:" + strings.Repeat("b", 64)
	return func(ctx context.Context, args []string) (string, error) {
		t.Helper()
		if ctx.Err() != nil {
			t.Fatal("command started after cancellation")
		}
		if args[0] == failure {
			return "private-error-with-secret", errors.New("private-error-with-secret")
		}
		switch args[0] {
		case "info":
			return `{"os":"linux","arch":"` + backend + `","cpus":4,"memory":1073741824,"version":"29.0.0","private":"` + syntheticHome("home") + `"}`, nil
		case "compose":
			switch args[1] {
			case "version":
				return "5.0.0\n", nil
			case "config":
				out, _ := json.Marshal(map[string]any{"services": map[string]any{"orders": map[string]any{"platform": requested, "environment": map[string]string{"PASSWORD": "private-secret"}, "volumes": []string{syntheticHome("Users") + ":/mnt"}}}})
				return string(out), nil
			case "ps":
				return id, nil
			}
		case "inspect":
			if args[len(args)-1] != id {
				t.Fatal("inspect used something other than observed container ID")
			}
			return `{"image":"` + image + `","memory":0,"nano_cpus":500000000,"running":true,"ignored":"private-secret"}`, nil
		case "image":
			if args[1] != "inspect" || args[len(args)-1] != image {
				t.Fatal("image metadata used a mutable tag")
			}
			return `[{"Id":"` + image + `","Os":"linux","Architecture":"` + actual + `","RepoDigests":["private.registry/app@` + image + `"],"Config":{"Env":["PASSWORD=private-secret"]}}]`, nil
		case "stats":
			return `{"cpu":"12.50%","memory":"1.5MiB / 1GiB","private":"` + syntheticHome("home") + `"}`, nil
		}
		t.Fatalf("unexpected metadata command %v", args)
		return "", nil
	}
}

func TestObservedImageIdentityPlatformsAndResources(t *testing.T) {
	t.Setenv("DOCKER_DEFAULT_PLATFORM", "")
	for _, test := range []struct{ name, backend, image, requested, wantBackend, wantRequest string }{
		{"matching", "x86_64", "amd64", "linux/amd64", "MATCH", "MATCH"},
		{"emulation candidate", "aarch64", "amd64", "linux/amd64", "MISMATCH", "MATCH"},
		{"request ignored", "amd64", "amd64", "linux/arm64", "MATCH", "MISMATCH"},
		{"no request", "amd64", "amd64", "", "MATCH", "NOT_REQUESTED"},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := Collect(t.Context(), dockerFixture(t, test.backend, test.image, test.requested, ""), []string{"compose"}, []string{"orders"})
			s := r.Services["orders"]
			if s.BackendRelation.Status != test.wantBackend || s.RequestedRelation.Status != test.wantRequest {
				t.Fatalf("relations: %+v", s)
			}
			if s.ContainerID == "" || s.ImageID == "" || len(s.ImageDigests) != 1 {
				t.Fatal("immutable provenance missing")
			}
			if r.BackendBudget.CPUs.Value == nil || *r.BackendBudget.CPUs.Value != 4 || s.Limits.MemoryBytes.Status != "NOT_CONFIGURED" || s.Limits.MemoryBytes.Value == nil || *s.Limits.MemoryBytes.Value != 0 {
				t.Fatal("budget and cap meanings lost")
			}
			if s.Usage.Status != "MEASURED" || s.Usage.ObservedAt == nil || *s.Usage.MemoryBytes != 1572864 || *s.Usage.CPUPercent != 12.5 {
				t.Fatalf("usage: %+v", s.Usage)
			}
			if r.Host.Status != "UNKNOWN" || r.ExecutionProfile.Status != "UNVERIFIED" || r.ComparisonEligibility.Status != "NOT_EVALUATED" {
				t.Fatal("platform facts were presented as execution proof")
			}
			data, _ := json.Marshal(r)
			for _, secret := range []string{"private-secret", "/home/", "/Users/", "private.registry", "PASSWORD", "private-error"} {
				if strings.Contains(string(data), secret) {
					t.Fatalf("metadata leaked %q", secret)
				}
			}
		})
	}
}

func TestDefaultPlatformAndMissingObservations(t *testing.T) {
	t.Setenv("DOCKER_DEFAULT_PLATFORM", "linux/arm64")
	r := Collect(t.Context(), dockerFixture(t, "arm64", "arm64", "", ""), []string{"compose"}, []string{"orders"})
	if r.Services["orders"].Requested.Source != "DOCKER_DEFAULT_PLATFORM" || r.Services["orders"].RequestedRelation.Status != "MATCH" {
		t.Fatal("default platform not observed")
	}
	for _, op := range []string{"info", "compose", "inspect", "image", "stats"} {
		t.Run(op, func(t *testing.T) {
			r := Collect(t.Context(), dockerFixture(t, "amd64", "amd64", "", op), []string{"compose"}, []string{"orders"})
			s := r.Services["orders"]
			if op == "info" && r.Backend.Status != "UNKNOWN" || op == "compose" && s.Requested.Status != "UNKNOWN" || (op == "inspect" || op == "image") && s.Image.Status != "UNKNOWN" || op == "stats" && s.Usage.Status != "NOT_MEASURED" {
				t.Fatalf("missing observation invented: %+v", r)
			}
			data, _ := json.Marshal(r)
			if strings.Contains(string(data), "private-error") {
				t.Fatal("raw error persisted")
			}
		})
	}
}

func TestPartialStartAndCancellation(t *testing.T) {
	base := dockerFixture(t, "amd64", "amd64", "", "")
	r := Collect(t.Context(), func(ctx context.Context, args []string) (string, error) {
		if args[0] == "compose" && args[1] == "ps" {
			return "", nil
		}
		return base(ctx, args)
	}, []string{"compose"}, []string{"orders"})
	if r.Backend.Status != "KNOWN" || r.Services["orders"].Image.Status != "UNKNOWN" || r.Services["orders"].Image.Reason == "" {
		t.Fatal("partial start lost known data or missing reason")
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	r = Collect(ctx, func(context.Context, []string) (string, error) {
		t.Fatal("cancelled collection executed command")
		return "", nil
	}, []string{"compose"}, []string{"orders"})
	if r.Services["orders"].Image.Status != "UNKNOWN" {
		t.Fatal("cancelled collection invented image")
	}
}

func TestMalformedOrMismatchedImageIsUnknown(t *testing.T) {
	for _, response := range []string{`{`, `[{"Id":"sha256:other","Os":"linux","Architecture":"amd64"}]`} {
		base := dockerFixture(t, "amd64", "amd64", "", "")
		r := Collect(t.Context(), func(ctx context.Context, args []string) (string, error) {
			if args[0] == "image" {
				return response, nil
			}
			return base(ctx, args)
		}, []string{"compose"}, []string{"orders"})
		if r.Services["orders"].Image.Status != "UNKNOWN" {
			t.Fatal("invalid image metadata accepted")
		}
	}
}

func TestPlatformNormalizationAndVariants(t *testing.T) {
	for _, test := range []struct{ value, arch, variant, state string }{
		{"linux/x86_64", "amd64", "", "KNOWN"}, {"linux/x86-64", "amd64", "", "KNOWN"}, {"linux/aarch64", "arm64", "", "KNOWN"}, {"linux/arm/v7", "arm", "v7", "KNOWN"},
		{"linux/arm64/v8", "arm64", "v8", "KNOWN"}, {"linux/arm64/private/path", "", "", "UNKNOWN"}, {syntheticHome("home"), "", "", "UNKNOWN"},
	} {
		p := parsePlatform(test.value, "fixture")
		if p.Status != test.state || p.Architecture != test.arch || p.Variant != test.variant {
			t.Fatalf("platform %q: %+v", test.value, p)
		}
	}
	if got := compare(parsePlatform("linux/arm/v7", "fixture"), parsePlatform("linux/arm", "fixture"), false); got.Status != "UNKNOWN" {
		t.Fatal("unknown variant guessed")
	}
	if got := compare(parsePlatform("linux/arm/v7", "fixture"), parsePlatform("linux/arm/v6", "fixture"), false); got.Status != "MISMATCH" {
		t.Fatal("variant mismatch hidden")
	}
	for _, value := range []string{"NaNMiB", "-1MiB", "1e100GiB", syntheticHome("home"), ""} {
		if _, ok := memoryBytes(value); ok {
			t.Fatalf("invalid measurement accepted: %q", value)
		}
	}
}

func TestImageVariantAndUnmeasuredVersusZero(t *testing.T) {
	base := dockerFixture(t, "arm64", "arm64", "linux/arm64", "")
	r := Collect(t.Context(), func(ctx context.Context, args []string) (string, error) {
		if args[0] == "image" {
			out, err := base(ctx, args)
			return strings.Replace(out, `"Architecture":"arm64"`, `"Architecture":"arm64","Variant":"v8"`, 1), err
		}
		if args[0] == "inspect" {
			return `{"image":"sha256:` + strings.Repeat("b", 64) + `","running":true}`, nil
		}
		if args[0] == "stats" {
			return `{"cpu":"0.00%","memory":"0B / 0B"}`, nil
		}
		return base(ctx, args)
	}, []string{"compose"}, []string{"orders"})
	s := r.Services["orders"]
	if s.Image.Variant != "v8" || s.Image.Raw != "linux/arm64/v8" || s.BackendRelation.Status != "UNKNOWN" || s.RequestedRelation.Status != "MATCH" {
		t.Fatalf("variant handling: %+v", s)
	}
	if s.Limits.MemoryBytes.Status != "UNKNOWN" || s.Limits.NanoCPUs.Status != "UNKNOWN" || s.Usage.Status != "MEASURED" || s.Usage.MemoryBytes == nil || *s.Usage.MemoryBytes != 0 {
		t.Fatal("unknown cap confused with measured zero")
	}
}
