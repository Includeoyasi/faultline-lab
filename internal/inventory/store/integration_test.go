package store

import (
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/database"
	"github.com/Includeoyasi/faultline-lab/internal/inventory/model"
)

var testDataset = func() string {
	if value := os.Getenv("FL_INVENTORY_DATASET_ID"); value != "" {
		return value
	}
	return "inventory-smoke-v1"
}()

func integrationStore(t *testing.T) *Store {
	t.Helper()
	if os.Getenv("FL_DATA_INTEGRATION") != "1" {
		t.Skip("requires isolated PostgreSQL smoke")
	}
	cfg, err := database.Load(os.Getenv("FL_DATA_CONFIG"), os.Getenv("FL_DATA_SECRETS"))
	if err != nil {
		t.Fatal(err)
	}
	identity, err := cfg.Service("inventory")
	if err != nil {
		t.Fatal(err)
	}
	store, err := Open(t.Context(), Config{Host: cfg.Host, Port: cfg.Port, Database: identity.Database, User: identity.AppUser, SSLMode: cfg.SSLMode, PasswordFile: filepath.Join(cfg.Secrets, identity.AppUser+".password"), DatasetID: testDataset, ConnectTimeout: "5s", RequestTimeout: "10s", MaxConnections: 4})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(store.Close)
	return store
}
func fixture(t *testing.T, s *Store, sku string, quantity, price int64) {
	t.Helper()
	if _, err := s.pool.Exec(t.Context(), "INSERT INTO inventory.products(sku,title,price_minor,currency) VALUES($1,$1,$2,'RUB');", sku, price); err != nil {
		t.Fatal(err)
	}
	if _, err := s.pool.Exec(t.Context(), "INSERT INTO inventory.stock(sku,available) VALUES($1,$2)", sku, quantity); err != nil {
		t.Fatal(err)
	}
}
func mutate(t *testing.T, s *Store, r model.Request) model.Receipt {
	t.Helper()
	result, err := s.Mutate(t.Context(), r)
	if err != nil {
		t.Fatal(err)
	}
	return result
}
func request(method, order, key string, items ...model.Item) model.Request {
	return model.Request{DatasetID: testDataset, OrderID: order, OperationID: key, Method: method, Items: items}
}
func scalar(t *testing.T, s *Store, query string, args ...any) int64 {
	t.Helper()
	var result int64
	if err := s.pool.QueryRow(t.Context(), query, args...).Scan(&result); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestInventoryFoundation(t *testing.T) {
	s := integrationStore(t)
	t.Run("commit absent is a durable refusal", func(t *testing.T) {
		r := request("commit", "absent", "absent-commit")
		first := mutate(t, s, r)
		if first.Outcome != model.InvalidState || first.Reservation != nil {
			t.Fatal(first)
		}
		if repeated := mutate(t, s, r); !reflect.DeepEqual(first, repeated) {
			t.Fatal("refusal not replayed")
		}
		if scalar(t, s, "SELECT count(*) FROM inventory.operations WHERE operation_id='absent-commit'") != 1 {
			t.Fatal("refusal was not stored")
		}
	})
	t.Run("atomic multi SKU refusal and durable rejection", func(t *testing.T) {
		fixture(t, s, "atomic-a", 2, 100)
		fixture(t, s, "atomic-b", 0, 200)
		r := request("reserve", "atomic", "atomic-key", model.Item{SKU: "atomic-a", Quantity: 1}, model.Item{SKU: "atomic-b", Quantity: 1})
		before := mutate(t, s, r)
		if before.Outcome != model.InsufficientStock {
			t.Fatal(before)
		}
		if scalar(t, s, "SELECT available FROM inventory.stock WHERE sku='atomic-a'") != 2 {
			t.Fatal("partial reservation")
		}
		if _, err := s.pool.Exec(t.Context(), "UPDATE inventory.stock SET available=1 WHERE sku='atomic-b'"); err != nil {
			t.Fatal(err)
		}
		if after := mutate(t, s, r); !reflect.DeepEqual(before, after) {
			t.Fatal("saved refusal changed")
		}
		r.OperationID = "atomic-new"
		if mutate(t, s, r).Outcome != model.Applied {
			t.Fatal("fresh retry failed")
		}
	})
	t.Run("concurrent limited stock and conservation", func(t *testing.T) {
		fixture(t, s, "contended", 10, 100)
		var successes atomic.Int64
		var wg sync.WaitGroup
		slots := make(chan struct{}, 8)
		for i := 0; i < 100; i++ {
			slots <- struct{}{}
			wg.Go(func() {
				defer func() { <-slots }()
				r := request("reserve", fmt.Sprintf("contended-%d", i), fmt.Sprintf("contended-key-%d", i), model.Item{SKU: "contended", Quantity: 1})
				result, err := s.Mutate(t.Context(), r)
				if err != nil {
					t.Error(err)
					return
				}
				if result.Outcome == model.Applied {
					successes.Add(1)
				} else if result.Outcome != model.InsufficientStock {
					t.Errorf("unexpected result %s", result.Outcome)
				}
			})
		}
		wg.Wait()
		if successes.Load() != 10 || scalar(t, s, "SELECT available FROM inventory.stock WHERE sku='contended'") != 0 || scalar(t, s, "SELECT sum(quantity) FROM inventory.reservation_items WHERE sku='contended'") != 10 {
			t.Fatal("stock conservation failed", successes.Load())
		}
	})
	t.Run("concurrent same key and changed payload", func(t *testing.T) {
		fixture(t, s, "dedupe", 5, 125)
		r := request("reserve", "dedupe", "dedupe-key", model.Item{SKU: "dedupe", Quantity: 2})
		var wg sync.WaitGroup
		for i := 0; i < 16; i++ {
			wg.Go(func() {
				result, err := s.Mutate(t.Context(), r)
				if err != nil || result.Outcome != model.Applied {
					t.Errorf("repeat=%v %v", result, err)
				}
			})
		}
		wg.Wait()
		if scalar(t, s, "SELECT available FROM inventory.stock WHERE sku='dedupe'") != 3 || scalar(t, s, "SELECT count(*) FROM inventory.ledger WHERE order_id='dedupe'") != 1 || scalar(t, s, "SELECT count(*) FROM inventory.outbox WHERE order_id='dedupe'") != 1 {
			t.Fatal("duplicate effect")
		}
		r.Items[0].Quantity = 1
		if _, err := s.Mutate(t.Context(), r); !errors.Is(err, model.ErrKeyConflict) {
			t.Fatal("changed payload accepted", err)
		}
		r.OperationID = "dedupe-second"
		if mutate(t, s, r).Outcome != model.ReservationExists {
			t.Fatal("second key created another reservation")
		}
	})
	t.Run("historical receipt and release barrier", func(t *testing.T) {
		fixture(t, s, "release", 3, 75)
		r := request("reserve", "release", "release-reserve", model.Item{SKU: "release", Quantity: 2})
		original := mutate(t, s, r)
		mutate(t, s, request("release", "release", "release-terminal"))
		mutate(t, s, request("release", "release", "release-terminal"))
		if repeat := mutate(t, s, r); !reflect.DeepEqual(original, repeat) {
			t.Fatal("historical receipt changed")
		}
		current, err := s.GetReservation(t.Context(), testDataset, "release")
		if err != nil || current.State != model.Released {
			t.Fatal(current, err)
		}
		if scalar(t, s, "SELECT available FROM inventory.stock WHERE sku='release'") != 3 {
			t.Fatal("release repeated effect")
		}
		if mutate(t, s, request("commit", "release", "release-commit")).Outcome != model.InvalidState {
			t.Fatal("committed released reservation")
		}
		mutate(t, s, request("release", "persist-barrier", "persist-release"))
		if mutate(t, s, request("reserve", "persist-barrier", "persist-late", model.Item{SKU: "release", Quantity: 1})).Outcome != model.ReservationExists {
			t.Fatal("late reserve bypassed barrier")
		}
	})
	t.Run("commit release race", func(t *testing.T) {
		fixture(t, s, "race", 1, 99)
		mutate(t, s, request("reserve", "race", "race-reserve", model.Item{SKU: "race", Quantity: 1}))
		var wg sync.WaitGroup
		for _, method := range []string{"commit", "release"} {
			wg.Go(func() {
				if _, err := s.Mutate(t.Context(), request(method, "race", "race-"+method)); err != nil {
					t.Error(err)
				}
			})
		}
		wg.Wait()
		current, err := s.GetReservation(t.Context(), testDataset, "race")
		if err != nil {
			t.Fatal(err)
		}
		want := int64(0)
		if current.State == model.Released {
			want = 1
		} else if current.State != model.Committed {
			t.Fatal(current)
		}
		if scalar(t, s, "SELECT available FROM inventory.stock WHERE sku='race'") != want || scalar(t, s, "SELECT count(*) FROM inventory.outbox WHERE order_id='race'") != 2 {
			t.Fatal("terminal race balance")
		}
	})
	t.Run("commit retries never sell twice or release sold stock", func(t *testing.T) {
		fixture(t, s, "committed", 3, 10)
		mutate(t, s, request("reserve", "committed", "committed-reserve", model.Item{SKU: "committed", Quantity: 2}))
		r := request("commit", "committed", "committed-first")
		first := mutate(t, s, r)
		if repeated := mutate(t, s, r); !reflect.DeepEqual(first, repeated) {
			t.Fatal("commit receipt changed")
		}
		if mutate(t, s, request("commit", "committed", "committed-second")).Outcome != model.Applied {
			t.Fatal("repeat terminal commit")
		}
		if mutate(t, s, request("release", "committed", "committed-release")).Outcome != model.InvalidState {
			t.Fatal("released sold stock")
		}
		if scalar(t, s, "SELECT available FROM inventory.stock WHERE sku='committed'") != 1 || scalar(t, s, "SELECT count(*) FROM inventory.ledger WHERE order_id='committed'") != 2 || scalar(t, s, "SELECT count(*) FROM inventory.outbox WHERE order_id='committed'") != 2 {
			t.Fatal("repeat commit changed business effect")
		}
	})
	t.Run("failed outbox insert rolls back the entire mutation", func(t *testing.T) {
		fixture(t, s, "rollback", 2, 1)
		mutate(t, s, request("release", "rollback-fixture", "rollback-fixture-key"))
		// An owned malformed event fixture forces failure AFTER stock and receipt writes.
		_, err := s.pool.Exec(t.Context(), "INSERT INTO inventory.outbox(event_id,order_id,aggregate_version,event_type,payload) VALUES('rollback-key','rollback-fixture',2,'InventoryReleased','{}')")
		if err != nil {
			t.Fatal(err)
		}
		if _, err = s.Mutate(t.Context(), request("reserve", "rollback", "rollback-key", model.Item{SKU: "rollback", Quantity: 1})); err == nil {
			t.Fatal("controlled outbox failure did not fail transaction")
		}
		if scalar(t, s, "SELECT available FROM inventory.stock WHERE sku='rollback'") != 2 || scalar(t, s, "SELECT count(*) FROM inventory.reservations WHERE order_id='rollback'") != 0 || scalar(t, s, "SELECT count(*) FROM inventory.operations WHERE operation_id='rollback-key'") != 0 || scalar(t, s, "SELECT count(*) FROM inventory.ledger WHERE order_id='rollback'") != 0 {
			t.Fatal("partial business transaction persisted")
		}
		if _, err = s.pool.Exec(t.Context(), "DELETE FROM inventory.outbox WHERE event_id='rollback-key'"); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("cancel while waiting on an owned database lock", func(t *testing.T) {
		lock, err := s.pool.Begin(t.Context())
		if err != nil {
			t.Fatal(err)
		}
		defer rollback(lock)
		if _, err = lock.Exec(t.Context(), "SELECT pg_advisory_xact_lock(2,hashtext('blocked'))"); err != nil {
			t.Fatal(err)
		}
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()
		done := make(chan error, 1)
		go func() { _, e := s.Mutate(ctx, request("release", "blocked", "blocked-key")); done <- e }()
		deadline := time.NewTimer(3 * time.Second)
		defer deadline.Stop()
		ticker := time.NewTicker(5 * time.Millisecond)
		defer ticker.Stop()
		waiting := false
		for !waiting {
			select {
			case <-deadline.C:
				t.Fatal("request never reached the contended lock")
			case <-ticker.C:
				waiting = scalar(t, s, "SELECT count(*) FROM pg_locks WHERE locktype='advisory' AND NOT granted AND classid=2 AND objid::bigint=(hashtext('blocked')::bigint & 4294967295)") > 0
			}
		}
		cancel()
		select {
		case err = <-done:
			if !errors.Is(err, context.Canceled) {
				t.Fatal("cancel result", err)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("blocked request ignored cancellation")
		}
		rollback(lock)
		if scalar(t, s, "SELECT count(*) FROM inventory.operations WHERE operation_id='blocked-key'") != 0 {
			t.Fatal("canceled operation persisted")
		}
		if mutate(t, s, request("release", "blocked", "blocked-retry")).Outcome != model.Applied {
			t.Fatal("pool did not recover after cancellation")
		}
	})
	t.Run("concurrent release before reserve", func(t *testing.T) {
		fixture(t, s, "barrier-race", 20, 1)
		for i := 0; i < 20; i++ {
			order := fmt.Sprintf("barrier-race-%d", i)
			var wg sync.WaitGroup
			wg.Go(func() {
				if _, err := s.Mutate(t.Context(), request("reserve", order, order+"-reserve", model.Item{SKU: "barrier-race", Quantity: 1})); err != nil {
					t.Error(err)
				}
			})
			wg.Go(func() {
				if _, err := s.Mutate(t.Context(), request("release", order, order+"-release")); err != nil {
					t.Error(err)
				}
			})
			wg.Wait()
			if mutate(t, s, request("reserve", order, order+"-late", model.Item{SKU: "barrier-race", Quantity: 1})).Outcome != model.ReservationExists {
				t.Fatal("late reserve bypassed cancellation")
			}
		}
		if scalar(t, s, "SELECT available FROM inventory.stock WHERE sku='barrier-race'") != 20 {
			t.Fatal("race leaked stock")
		}
	})
	t.Run("overflow unknown SKU cancellation and dataset isolation", func(t *testing.T) {
		fixture(t, s, "overflow", math.MaxInt64, math.MaxInt64)
		if _, err := s.Mutate(t.Context(), request("reserve", "overflow", "overflow-key", model.Item{SKU: "overflow", Quantity: 2})); !errors.Is(err, model.ErrInvalid) {
			t.Fatal("overflow", err)
		}
		if mutate(t, s, request("reserve", "unknown", "unknown-key", model.Item{SKU: "missing", Quantity: 1})).Outcome != model.UnknownSKU {
			t.Fatal("unknown SKU")
		}
		canceled, cancel := context.WithCancel(t.Context())
		cancel()
		if _, err := s.Mutate(canceled, request("release", "canceled", "canceled-key")); !errors.Is(err, context.Canceled) {
			t.Fatal("cancel", err)
		}
		r := request("release", "foreign", "foreign-key")
		r.DatasetID = "old-dataset"
		if _, err := s.Mutate(t.Context(), r); !errors.Is(err, model.ErrDataset) {
			t.Fatal("foreign dataset", err)
		}
		if scalar(t, s, "SELECT count(*) FROM inventory.operations WHERE operation_id IN ('overflow-key','canceled-key','foreign-key')") != 0 {
			t.Fatal("invalid operation persisted")
		}
	})
}

func TestInventoryAfterRestart(t *testing.T) {
	s := integrationStore(t)
	r, err := s.GetReservation(t.Context(), testDataset, "persist-barrier")
	if err != nil || r.State != model.Released {
		t.Fatal("barrier lost", r, err)
	}
	if result := mutate(t, s, request("release", "persist-barrier", "persist-release")); result.Reservation.State != model.Released {
		t.Fatal("receipt lost")
	}
	if mutate(t, s, request("reserve", "persist-barrier", "after-restart", model.Item{SKU: "release", Quantity: 1})).Outcome != model.ReservationExists {
		t.Fatal("restart bypassed barrier")
	}
	if scalar(t, s, "SELECT count(*) FROM inventory.ledger WHERE order_id='dedupe'") != 1 {
		t.Fatal("ledger lost")
	}
	if scalar(t, s, "SELECT count(*) FROM inventory.outbox WHERE order_id='dedupe'") != 1 {
		t.Fatal("event lost")
	}
	// A new configured namespace cannot attach to surviving data accidentally.
	if _, err = NewStore(t.Context(), s.pool, "another-dataset", time.Second); !errors.Is(err, model.ErrDataset) {
		t.Fatal("dataset rebinding allowed", err)
	}
}
