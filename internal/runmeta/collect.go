// Package runmeta records observed Docker facts without changing smoke outcomes.
package runmeta

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Platform separates a value from its origin and any missing observation.
type Platform struct {
	OS           string `json:"os,omitempty"`
	Architecture string `json:"architecture,omitempty"`
	Variant      string `json:"variant,omitempty"`
	Raw          string `json:"raw,omitempty"`
	Status       string `json:"status"`
	Source       string `json:"source"`
	Reason       string `json:"reason,omitempty"`
}

// Quantity uses null for unknown. Zero means no cap in this particular field;
// e.g. NanoCpus=0 says nothing about quota/cpuset or a parent cgroup.
type Quantity struct {
	Value  *int64 `json:"value"`
	Status string `json:"status"`
	Source string `json:"source"`
	Reason string `json:"reason,omitempty"`
}

type Relation struct {
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type Budget struct {
	CPUs        Quantity `json:"cpus"`
	MemoryBytes Quantity `json:"memory_bytes"`
}

type Limits struct {
	NanoCPUs    Quantity `json:"nano_cpus"`
	MemoryBytes Quantity `json:"memory_bytes"`
}

// Sample is one Docker stats observation, not a load/capacity measurement.
type Sample struct {
	ObservedAt  *time.Time `json:"observed_at,omitempty"`
	CPUPercent  *float64   `json:"cpu_percent"`
	MemoryBytes *int64     `json:"memory_bytes"`
	Source      string     `json:"source"`
	Status      string     `json:"status"`
	Reason      string     `json:"reason,omitempty"`
}

type Service struct {
	ContainerID       string   `json:"container_id,omitempty"`
	ImageID           string   `json:"image_id,omitempty"`
	ImageDigests      []string `json:"image_digests"`
	Requested         Platform `json:"requested_platform"`
	Image             Platform `json:"image_platform"`
	BackendRelation   Relation `json:"backend_image_relation"`
	RequestedRelation Relation `json:"requested_actual_relation"`
	Limits            Limits   `json:"limits"`
	Usage             Sample   `json:"usage_snapshot"`
}

// Report version 1 covers the named long-lived tested containers only.
// Tool containers removed by compose run --rm are deliberately not inferred.
type Report struct {
	SchemaVersion         int                `json:"schema_version"`
	Runner                Platform           `json:"runner_platform"`
	Host                  Platform           `json:"host_platform"`
	Backend               Platform           `json:"backend_platform"`
	DockerVersion         string             `json:"docker_version,omitempty"`
	ComposeVersion        string             `json:"compose_version,omitempty"`
	GoVersion             string             `json:"go_version"`
	BackendBudget         Budget             `json:"backend_budget"`
	Services              map[string]Service `json:"services"`
	ExecutionProfile      Relation           `json:"execution_profile"`
	ComparisonEligibility Relation           `json:"comparison_eligibility"`
	Coverage              string             `json:"coverage"`
}

var (
	platformToken = regexp.MustCompile(`^[a-z0-9_-]{1,32}$`)
	variantToken  = regexp.MustCompile(`^v[0-9]{1,3}$`)
	containerID   = regexp.MustCompile(`^[a-f0-9]{12,64}$`)
	imageID       = regexp.MustCompile(`^sha256:[a-f0-9]{64}$`)
	versionToken  = regexp.MustCompile(`^[a-zA-Z0-9.+_-]{1,80}$`)
)

func unknownPlatform(source, reason string) Platform {
	return Platform{Status: "UNKNOWN", Source: source, Reason: reason}
}

func unknownQuantity(source, reason string) Quantity {
	return Quantity{Status: "UNKNOWN", Source: source, Reason: reason}
}

// New supplies honest initial values even if build/start is cancelled.
func New(names []string) Report {
	r := Report{
		SchemaVersion: 1, Runner: parsePlatform(runtime.GOOS+"/"+runtime.GOARCH, "runtime.GOOS/GOARCH"),
		Host:    unknownPlatform("host profile", "physical_host_not_observed"),
		Backend: unknownPlatform("docker info", "not_collected"), GoVersion: runtime.Version(),
		BackendBudget:         Budget{unknownQuantity("docker info.NCPU", "not_collected"), unknownQuantity("docker info.MemTotal", "not_collected")},
		Services:              map[string]Service{},
		ExecutionProfile:      Relation{"UNVERIFIED", "platform_match_does_not_prove_native_execution"},
		ComparisonEligibility: Relation{"NOT_EVALUATED", "requires_verified_execution_profile_and_experiment_contract"},
		Coverage:              "named long-lived tested containers; ephemeral tool containers are not observed",
	}
	for _, name := range names {
		r.Services[name] = Service{
			ImageDigests: []string{}, Requested: unknownPlatform("effective Compose config", "not_collected"),
			Image:           unknownPlatform("docker image inspect by container image ID", "container_not_observed"),
			BackendRelation: Relation{"UNKNOWN", "platform_not_observed"}, RequestedRelation: Relation{"UNKNOWN", "request_not_observed"},
			Limits: Limits{unknownQuantity("container.HostConfig.NanoCpus", "container_not_observed"), unknownQuantity("container.HostConfig.Memory", "container_not_observed")},
			Usage:  Sample{Source: "docker stats", Status: "NOT_MEASURED", Reason: "not_collected"},
		}
	}
	return r
}

// Docker executes an argv without persisting raw stdout/stderr. Metadata queries
// must bypass command loggers because effective Compose config contains secrets.
type Docker func(context.Context, []string) (string, error)

// Collect has one total budget and per-command deadlines. Missing metadata is
// data, not a functional failure; callers still own cancellation and cleanup.
func Collect(parent context.Context, docker Docker, compose []string, names []string) Report {
	r := New(names)
	ctx, cancel := context.WithTimeout(parent, 30*time.Second)
	defer cancel()
	query := func(args ...string) (string, bool) {
		if ctx.Err() != nil {
			return "", false
		}
		bounded, stop := context.WithTimeout(ctx, 5*time.Second)
		defer stop()
		out, err := docker(bounded, args)
		return out, err == nil && bounded.Err() == nil
	}
	var backend struct {
		OS      string `json:"os"`
		Arch    string `json:"arch"`
		CPUs    *int64 `json:"cpus"`
		Memory  *int64 `json:"memory"`
		Version string `json:"version"`
	}
	const backendFormat = `{"os":{{json .OSType}},"arch":{{json .Architecture}},"cpus":{{json .NCPU}},"memory":{{json .MemTotal}},"version":{{json .ServerVersion}}}`
	if out, ok := query("info", "--format", backendFormat); ok && json.Unmarshal([]byte(out), &backend) == nil {
		r.Backend = parsePlatform(backend.OS+"/"+backend.Arch, "docker info")
		r.BackendBudget = Budget{quantity(backend.CPUs, "docker info.NCPU", false), quantity(backend.Memory, "docker info.MemTotal", false)}
		if versionToken.MatchString(backend.Version) {
			r.DockerVersion = backend.Version
		}
	} else {
		r.Backend.Reason = "query_failed_or_invalid_response"
	}
	if out, ok := query("compose", "version", "--short"); ok && versionToken.MatchString(strings.TrimSpace(out)) {
		r.ComposeVersion = strings.TrimSpace(out)
	}
	var config struct {
		Services map[string]struct {
			Platform string `json:"platform"`
		} `json:"services"`
	}
	configOK := false
	if out, ok := query(append(append([]string{}, compose...), "config", "--format", "json")...); ok {
		configOK = json.Unmarshal([]byte(out), &config) == nil && config.Services != nil
	}
	for _, name := range names {
		s := r.Services[name]
		if svc, exists := config.Services[name]; configOK && exists {
			switch {
			case svc.Platform != "":
				s.Requested = parsePlatform(svc.Platform, "effective Compose service.platform")
			case os.Getenv("DOCKER_DEFAULT_PLATFORM") != "":
				s.Requested = parsePlatform(os.Getenv("DOCKER_DEFAULT_PLATFORM"), "DOCKER_DEFAULT_PLATFORM")
			default:
				s.Requested = Platform{Status: "UNSET", Source: "effective Compose config and DOCKER_DEFAULT_PLATFORM"}
			}
		} else {
			s.Requested.Reason = "config_unavailable_or_service_missing"
		}
		s.RequestedRelation = compare(s.Requested, s.Image, true)
		ids, ok := query(append(append([]string{}, compose...), "ps", "--all", "-q", name)...)
		id := strings.TrimSpace(ids)
		if !ok || !containerID.MatchString(id) {
			s.Image.Reason = "unique_container_not_observed"
			r.Services[name] = s
			continue
		}
		s.ContainerID = id
		var c struct {
			Image    string `json:"image"`
			Memory   *int64 `json:"memory"`
			NanoCPUs *int64 `json:"nano_cpus"`
			Running  *bool  `json:"running"`
		}
		const containerFormat = `{"image":{{json .Image}},"memory":{{json .HostConfig.Memory}},"nano_cpus":{{json .HostConfig.NanoCpus}},"running":{{json .State.Running}}}`
		if out, ok := query("inspect", "--format", containerFormat, id); ok && json.Unmarshal([]byte(out), &c) == nil {
			s.Limits = Limits{quantity(c.NanoCPUs, "container.HostConfig.NanoCpus", true), quantity(c.Memory, "container.HostConfig.Memory", true)}
			if imageID.MatchString(c.Image) {
				s.ImageID = c.Image
				var images []struct {
					ID      string   `json:"id"`
					OS      string   `json:"os"`
					Arch    string   `json:"Architecture"`
					Variant string   `json:"variant"`
					Digests []string `json:"RepoDigests"`
				}
				// Decode image JSON in memory with an allowlist; an absent optional
				// Variant must not fail a Go template query.
				if out, ok := query("image", "inspect", c.Image); ok && json.Unmarshal([]byte(out), &images) == nil && len(images) == 1 && images[0].ID == c.Image {
					img := images[0]
					value := img.OS + "/" + img.Arch
					if img.Variant != "" {
						value += "/" + img.Variant
					}
					s.Image = parsePlatform(value, "docker image inspect by container image ID")
					for _, d := range img.Digests {
						_, digest, _ := strings.Cut(d, "@")
						if imageID.MatchString(digest) {
							s.ImageDigests = append(s.ImageDigests, digest)
						}
					}
				} else {
					s.Image.Reason = "image_inspect_failed_or_id_mismatch"
				}
			} else {
				s.Image.Reason = "invalid_or_missing_image_id"
			}
		} else {
			s.Image.Reason = "container_inspect_failed_or_invalid"
		}
		s.BackendRelation = compare(r.Backend, s.Image, false)
		s.RequestedRelation = compare(s.Requested, s.Image, true)
		if c.Running != nil && *c.Running {
			s.Usage = sample(query, id)
		} else {
			s.Usage.Reason = "container_not_observed_running"
		}
		r.Services[name] = s
	}
	return r
}

func parsePlatform(value, source string) Platform {
	p := unknownPlatform(source, "invalid_or_missing_platform")
	parts := strings.Split(strings.ToLower(strings.TrimSpace(value)), "/")
	if len(parts) < 2 || len(parts) > 3 || !platformToken.MatchString(parts[0]) || !platformToken.MatchString(parts[1]) {
		return p
	}
	p.OS, p.Architecture = parts[0], parts[1]
	switch p.Architecture {
	case "x86_64", "x86-64":
		p.Architecture = "amd64"
	case "aarch64":
		p.Architecture = "arm64"
	}
	if len(parts) == 3 {
		if !variantToken.MatchString(parts[2]) {
			return unknownPlatform(source, "invalid_variant")
		}
		p.Variant = parts[2]
	}
	p.Raw = strings.TrimSpace(value)
	p.Status = "KNOWN"
	p.Reason = ""
	return p
}

func quantity(value *int64, source string, limit bool) Quantity {
	if value == nil || *value < 0 || (*value == 0 && !limit) {
		return unknownQuantity(source, "missing_or_invalid_value")
	}
	state := "KNOWN"
	if *value == 0 && limit {
		state = "NOT_CONFIGURED"
	}
	return Quantity{Value: value, Status: state, Source: source}
}

func compare(left, right Platform, requested bool) Relation {
	if requested && left.Status == "UNSET" {
		return Relation{Status: "NOT_REQUESTED"}
	}
	if left.Status != "KNOWN" || right.Status != "KNOWN" {
		return Relation{"UNKNOWN", "platform_not_observed"}
	}
	if left.OS != right.OS || left.Architecture != right.Architecture || (left.Variant != "" && right.Variant != "" && left.Variant != right.Variant) {
		return Relation{Status: "MISMATCH"}
	}
	if (!requested || left.Variant != "") && left.Variant != right.Variant {
		return Relation{"UNKNOWN", "variant_not_comparable"}
	}
	return Relation{Status: "MATCH"}
}

func sample(query func(...string) (string, bool), id string) Sample {
	s := Sample{Source: "docker stats", Status: "NOT_MEASURED", Reason: "query_failed_or_invalid_response"}
	const format = `{"cpu":{{json .CPUPerc}},"memory":{{json .MemUsage}}}`
	var data struct {
		CPU    string `json:"cpu"`
		Memory string `json:"memory"`
	}
	out, ok := query("stats", "--no-stream", "--format", format, id)
	if !ok || json.Unmarshal([]byte(out), &data) != nil {
		return s
	}
	cpu, err := strconv.ParseFloat(strings.TrimSuffix(data.CPU, "%"), 64)
	memory, _, _ := strings.Cut(data.Memory, "/")
	bytes, valid := memoryBytes(strings.TrimSpace(memory))
	if err != nil || !strings.HasSuffix(data.CPU, "%") || math.IsNaN(cpu) || math.IsInf(cpu, 0) || cpu < 0 || !valid {
		return s
	}
	now := time.Now().UTC()
	s.ObservedAt = &now
	s.CPUPercent = &cpu
	s.MemoryBytes = &bytes
	s.Status = "MEASURED"
	s.Reason = ""
	return s
}

func memoryBytes(value string) (int64, bool) {
	for _, unit := range []struct {
		name string
		size float64
	}{{"TiB", 1 << 40}, {"GiB", 1 << 30}, {"MiB", 1 << 20}, {"KiB", 1 << 10}, {"TB", 1e12}, {"GB", 1e9}, {"MB", 1e6}, {"kB", 1e3}, {"B", 1}} {
		if !strings.HasSuffix(value, unit.name) {
			continue
		}
		n, err := strconv.ParseFloat(strings.TrimSpace(strings.TrimSuffix(value, unit.name)), 64)
		n *= unit.size
		if err != nil || math.IsNaN(n) || math.IsInf(n, 0) || n < 0 || n >= math.MaxInt64 {
			return 0, false
		}
		return int64(n), true
	}
	return 0, false
}
