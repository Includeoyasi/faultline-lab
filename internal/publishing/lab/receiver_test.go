package lab

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"testing"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/Includeoyasi/faultline-lab/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
)

// Runs in a new disposable Compose project with no host ports. The reader
// connects as a separate role that has access only to its own schema.
func TestPublicationAndReplay(t *testing.T) {
	dsn := os.Getenv("FL_F09_TEST_DB")
	if dsn == "" {
		t.Skip("requires isolated F09 database")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 90*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	exec := func(query string) {
		t.Helper()
		if _, err := pool.Exec(ctx, query); err != nil {
			t.Fatal(err)
		}
	}
	exec("CREATE ROLE f09_reader LOGIN")
	exec(Schema)
	exec("GRANT USAGE ON SCHEMA f09_receiver TO f09_reader; GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA f09_receiver TO f09_reader")
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		t.Fatal(err)
	}
	cfg.ConnConfig.User = "f09_reader"
	reader, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer reader.Close()
	broker := os.Getenv("FL_F09_TEST_BROKER")
	for _, owner := range []string{"orders", "inventory"} {
		paths, err := fs.Glob(migrations.Files, owner+"/*.sql")
		if err != nil || len(paths) < 2 {
			t.Fatal(paths, err)
		}
		apply := func(path string) {
			b, err := migrations.Files.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			exec(string(b))
		}
		for _, path := range paths[:len(paths)-1] {
			apply(path)
		}
		exec("INSERT INTO " + owner + ".dataset(dataset_id) VALUES ('fixture')")
		kind := "InventoryReserved"
		if owner == "orders" {
			exec("INSERT INTO orders.orders(order_id,request_key,request_hash,reserve_operation_id,state,version,items) VALUES ('o','r',repeat('a',64),'reserve','ACCEPTED',1,'[{}]')")
			kind = "OrderAwaitingPayment"
		} else {
			exec("INSERT INTO inventory.reservations(order_id,state,version) VALUES('o','RESERVED',1)")
		}
		exec("INSERT INTO " + owner + ".outbox(event_id,order_id,aggregate_version,event_type,payload) VALUES('one','o',3,'" + kind + "','{}')")
		apply(paths[len(paths)-1])
		insert := "INSERT INTO " + owner + ".outbox(event_id,order_id,aggregate_version,event_type,payload) VALUES('two','o',6,'" + kind + "','{}')"
		exec("BEGIN; " + insert + "; ROLLBACK")
		exec(insert)
		p := publishing.Publisher{Pool: pool, Owner: owner, Dataset: "fixture", Broker: broker}
		if err = p.Once(ctx); err != nil {
			t.Fatal(err)
		}
		// Simulate a lost SQL mark after a committed Kafka delivery.
		exec("UPDATE " + owner + ".outbox SET published_at=NULL WHERE event_id='one'")
		for range 2 {
			if err = p.Once(ctx); err != nil {
				t.Fatal(err)
			}
		}
		cl, err := kgo.NewClient(kgo.SeedBrokers(broker), kgo.ConsumeTopics(owner+".events.v1"), kgo.FetchIsolationLevel(kgo.ReadCommitted()), kgo.ConsumeResetOffset(kgo.NewOffset().AtStart()))
		if err != nil {
			t.Fatal(err)
		}
		func() {
			defer cl.Close()
			seen := 0
			for seen < 3 {
				fetch := cl.PollRecords(ctx, 1)
				if es := fetch.Errors(); len(es) > 0 {
					t.Fatal(es[0].Err)
				}
				for _, r := range fetch.Records() {
					if string(r.Key) != "o" {
						continue
					}
					if err := Save(ctx, reader, "fixture", r); err != nil {
						t.Fatal(err)
					}
					if err := Apply(ctx, reader, "fixture", owner, "o"); err != nil && !errors.Is(err, pgx.ErrNoRows) {
						t.Fatal(err)
					}
					seen++
				}
			}
		}()
		var cursor, effects int
		if err := reader.QueryRow(ctx, "SELECT cursor,effects FROM f09_receiver.projection WHERE owner=$1", owner).Scan(&cursor, &effects); err != nil {
			t.Fatal(err)
		}
		expected := 2
		if owner == "orders" {
			expected = 0
		}
		if cursor != 2 || effects != expected {
			t.Fatalf("%s cursor %d effects %d", owner, cursor, effects)
		}
	}
}
