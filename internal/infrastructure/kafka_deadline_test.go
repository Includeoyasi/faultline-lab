package infrastructure

import (
	"context"
	"encoding/binary"
	"io"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/kmsg"
	"github.com/twmb/franz-go/pkg/kversion"
)

// This fixture completes metadata and producer initialization, accepts Produce,
// then loses its acknowledgement. A connection-only blackhole cannot test this.
func TestKafkaProduceCancellationAfterWrite(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	host, portText, _ := net.SplitHostPort(listener.Addr().String())
	port, _ := strconv.Atoi(portText)
	var connections []net.Conn
	var mu sync.Mutex
	var workers sync.WaitGroup
	accepted := make(chan struct{})
	acceptLoopDone := make(chan struct{})
	var acceptOnce sync.Once
	workers.Go(func() {
		defer close(acceptLoopDone)
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			connections = append(connections, conn)
			mu.Unlock()
			workers.Go(func() {
				defer func() { _ = conn.Close() }()
				for {
					var size [4]byte
					if _, err := io.ReadFull(conn, size[:]); err != nil {
						return
					}
					length := binary.BigEndian.Uint32(size[:])
					if length < 10 || length > 1<<20 {
						return
					}
					request := make([]byte, length)
					if _, err := io.ReadFull(conn, request); err != nil {
						return
					}
					key := binary.BigEndian.Uint16(request)
					version := int16(binary.BigEndian.Uint16(request[2:]))
					var response kmsg.Response
					switch key {
					case 0: // Produce was written; deliberately withhold the response.
						acceptOnce.Do(func() { close(accepted) })
						continue
					case 18:
						api := kmsg.NewPtrApiVersionsResponse()
						for _, pair := range [][2]int16{{0, 3}, {3, 4}, {18, 1}, {22, 0}} {
							entry := kmsg.NewApiVersionsResponseApiKey()
							entry.ApiKey, entry.MaxVersion = pair[0], pair[1]
							api.ApiKeys = append(api.ApiKeys, entry)
						}
						response = api
					case 3:
						metadata := kmsg.NewPtrMetadataResponse()
						broker := kmsg.NewMetadataResponseBroker()
						broker.NodeID, broker.Host, broker.Port = 1, host, int32(port)
						metadata.Brokers, metadata.ControllerID = []kmsg.MetadataResponseBroker{broker}, 1
						name := "deadline"
						topic := kmsg.NewMetadataResponseTopic()
						topic.Topic = &name
						partition := kmsg.NewMetadataResponseTopicPartition()
						partition.Leader, partition.Replicas, partition.ISR = 1, []int32{1}, []int32{1}
						topic.Partitions = []kmsg.MetadataResponseTopicPartition{partition}
						metadata.Topics = []kmsg.MetadataResponseTopic{topic}
						response = metadata
					case 22:
						producer := kmsg.NewPtrInitProducerIDResponse()
						producer.ProducerID = 1
						response = producer
					default:
						t.Errorf("unexpected fixture request key %d", key)
						return
					}
					response.SetVersion(version)
					body := response.AppendTo(append([]byte(nil), request[4:8]...))
					binary.BigEndian.PutUint32(size[:], uint32(len(body)))
					if _, err := conn.Write(append(size[:], body...)); err != nil {
						return
					}
				}
			})
		}
	})
	t.Cleanup(func() {
		_ = listener.Close()
		<-acceptLoopDone
		mu.Lock()
		for _, conn := range connections {
			_ = conn.Close()
		}
		mu.Unlock()
		workers.Wait()
	})
	cfg := testConfig(t)
	cfg.KafkaBrokers = []string{listener.Addr().String()}
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	client, err := kafkaClient(ctx, cfg, kgo.MaxVersions(kversion.V0_11_0()))
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	done := make(chan error, 1)
	go func() {
		done <- client.ProduceSync(ctx, &kgo.Record{Topic: "deadline", Value: []byte("probe")}).FirstErr()
	}()
	select {
	case <-accepted:
	case <-time.After(3 * time.Second):
		t.Fatal("fixture never received Produce")
	}
	cancel()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("lost acknowledgement was reported as success")
		}
	case <-time.After(time.Second):
		t.Fatal("in-flight Produce ignored cancellation")
	}
}
