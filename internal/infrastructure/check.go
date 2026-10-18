package infrastructure

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Check verifies Redis ACL/TTL and creates or re-reads a Kafka restart checkpoint.
func Check(ctx context.Context, cfg Config, stage, saved string) (Checkpoint, error) {
	var checkpoint Checkpoint
	if err := cfg.Validate(); err != nil {
		return checkpoint, err
	}
	if stage != "before-restart" && stage != "after-restart" {
		return checkpoint, errors.New("infrastructure stage must be before-restart or after-restart")
	}
	if stage == "after-restart" {
		if len(saved) > 2048 {
			return checkpoint, errors.New("Kafka checkpoint is too large")
		}
		decoder := json.NewDecoder(strings.NewReader(saved))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&checkpoint); err != nil {
			return checkpoint, errors.New("invalid Kafka checkpoint")
		}
		var extra any
		if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
			return checkpoint, errors.New("Kafka checkpoint has trailing data")
		}
		if checkpoint.Version != 1 || checkpoint.Topic != cfg.KafkaTopic || checkpoint.Partition != 0 || checkpoint.Offset < 0 || checkpoint.Key == "" || len(checkpoint.Key) > 64 || checkpoint.Value == "" || len(checkpoint.Value) > 256 {
			return checkpoint, errors.New("Kafka checkpoint does not match the configured probe")
		}
	} else if saved != "" {
		return checkpoint, errors.New("before-restart must create a new checkpoint")
	}
	if err := checkRedis(ctx, cfg); err != nil {
		return checkpoint, err
	}
	ctx, cancel := context.WithTimeout(ctx, cfg.budget())
	defer cancel()
	if stage == "before-restart" {
		return prepareKafka(ctx, cfg)
	}
	return checkpoint, verifyKafka(ctx, cfg, checkpoint)
}

func safeError(stage string, err error) error {
	for _, sentinel := range []error{context.Canceled, context.DeadlineExceeded} {
		if errors.Is(err, sentinel) {
			return fmt.Errorf("%s: %w", stage, sentinel)
		}
	}
	return errors.New(stage + " failed")
}
