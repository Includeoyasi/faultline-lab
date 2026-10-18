package infrastructure

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
)

// Checkpoint identifies the acknowledged record that must survive a restart.
type Checkpoint struct {
	Version   int    `json:"version"`
	Topic     string `json:"topic"`
	Partition int32  `json:"partition"`
	Offset    int64  `json:"offset"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}

func kafkaClient(ctx context.Context, cfg Config, options ...kgo.Opt) (*kgo.Client, error) {
	base := []kgo.Opt{
		kgo.WithContext(ctx), kgo.DisableClientMetrics(),
		kgo.SeedBrokers(cfg.KafkaBrokers...), kgo.DialTimeout(cfg.budget()),
		kgo.RequestTimeoutOverhead(time.Second), kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordDeliveryTimeout(cfg.budget()), kgo.MaxProduceRequestsInflightPerBroker(1),
		// A lost acknowledgement must fail this disposable probe within its budget.
		// It is never retried by the application; this is not an outbox policy.
		kgo.AllowIdempotentProduceCancellation(),
		kgo.ProducerBatchCompression(kgo.NoCompression()),
		kgo.FetchMaxWait(100 * time.Millisecond), kgo.FetchMaxBytes(1 << 20),
	}
	return kgo.NewClient(append(base, options...)...)
}

func prepareKafka(ctx context.Context, cfg Config) (Checkpoint, error) {
	var checkpoint Checkpoint
	client, err := kafkaClient(ctx, cfg)
	if err != nil {
		return checkpoint, safeError("kafka client configuration", err)
	}
	defer client.Close()
	request := kmsg.NewPtrCreateTopicsRequest()
	request.TimeoutMillis = int32(cfg.budget().Milliseconds())
	topic := kmsg.NewCreateTopicsRequestTopic()
	topic.Topic, topic.NumPartitions, topic.ReplicationFactor = cfg.KafkaTopic, 1, 1
	request.Topics = []kmsg.CreateTopicsRequestTopic{topic}
	response, err := client.Request(ctx, request)
	if err != nil {
		return checkpoint, safeError("kafka topic creation", err)
	}
	created, ok := response.(*kmsg.CreateTopicsResponse)
	if !ok || len(created.Topics) != 1 || created.Topics[0].ErrorCode != 0 {
		return checkpoint, errors.New("kafka check requires a newly created isolated topic")
	}
	var marker [16]byte
	if _, err = rand.Read(marker[:]); err != nil {
		return checkpoint, errors.New("cannot generate Kafka checkpoint")
	}
	record := &kgo.Record{Topic: cfg.KafkaTopic, Key: []byte("checkpoint"), Value: []byte(hex.EncodeToString(marker[:]))}
	if err = client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return checkpoint, safeError("kafka produce", err)
	}
	checkpoint = Checkpoint{1, record.Topic, record.Partition, record.Offset, string(record.Key), string(record.Value)}
	return checkpoint, readKafka(ctx, cfg, checkpoint)
}

func readKafka(ctx context.Context, cfg Config, expected Checkpoint) error {
	client, err := kafkaClient(ctx, cfg, kgo.ConsumePartitions(map[string]map[int32]kgo.Offset{
		expected.Topic: {expected.Partition: kgo.NewOffset().At(expected.Offset)},
	}))
	if err != nil {
		return safeError("kafka consumer configuration", err)
	}
	defer client.Close()
	for ctx.Err() == nil {
		fetches := client.PollFetches(ctx)
		if len(fetches.Errors()) != 0 {
			return safeError("kafka fetch", ctx.Err())
		}
		for _, record := range fetches.Records() {
			if record.Topic != expected.Topic || record.Partition != expected.Partition || record.Offset != expected.Offset || string(record.Key) != expected.Key || string(record.Value) != expected.Value {
				return errors.New("kafka checkpoint record does not match")
			}
			return nil
		}
	}
	return safeError("kafka fetch", ctx.Err())
}

func verifyKafka(ctx context.Context, cfg Config, expected Checkpoint) error {
	if err := readKafka(ctx, cfg, expected); err != nil {
		return err
	}
	client, err := kafkaClient(ctx, cfg)
	if err != nil {
		return safeError("kafka producer configuration", err)
	}
	defer client.Close()
	record := &kgo.Record{Topic: expected.Topic, Key: []byte(expected.Key), Value: []byte(expected.Value + ":after")}
	if err = client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return safeError("kafka produce after restart", err)
	}
	if record.Partition != expected.Partition || record.Offset <= expected.Offset {
		return errors.New("kafka new record did not follow the persisted checkpoint")
	}
	return readKafka(ctx, cfg, Checkpoint{1, record.Topic, record.Partition, record.Offset, string(record.Key), string(record.Value)})
}
