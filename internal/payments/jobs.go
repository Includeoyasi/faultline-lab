// Package payments processes reserved orders against the simulated payment ledger.
package payments

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"regexp"
	"time"

	"github.com/Includeoyasi/faultline-lab/internal/publishing"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twmb/franz-go/pkg/kgo"
)

var identifier = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}$`)
var currencyCode = regexp.MustCompile(`^[A-Z]{3}$`)

type processor struct {
	pool          *pgxpool.Pool
	dataset, seed string
	declineBPS    int
	maxDelay      time.Duration
}

func (p processor) bind(ctx context.Context) error {
	return pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "INSERT INTO payments.dataset(dataset_id) VALUES($1) ON CONFLICT DO NOTHING", p.dataset); err != nil {
			return err
		}
		var dataset string
		if err := tx.QueryRow(ctx, "SELECT dataset_id FROM payments.dataset WHERE singleton").Scan(&dataset); err != nil {
			return err
		}
		if dataset != p.dataset {
			return errors.New("payments database belongs to another dataset")
		}
		return nil
	})
}

func (p processor) simulate(order string) (bool, time.Duration) {
	hash := sha256.Sum256([]byte(p.seed + "\x00" + order))
	approved := int(binary.BigEndian.Uint32(hash[:4])%10000) >= p.declineBPS
	delay := time.Duration(uint64(binary.BigEndian.Uint32(hash[4:8]))%(uint64(p.maxDelay.Milliseconds())+1)) * time.Millisecond
	return approved, delay
}

// save commits the inbox and job before the caller acknowledges the Kafka record.
func (p processor) save(ctx context.Context, r *kgo.Record) error {
	var e publishing.Event
	if err := json.Unmarshal(r.Value, &e); err != nil {
		return errors.New("invalid Orders event JSON")
	}
	if e.Owner != "orders" || e.Dataset != p.dataset || e.SchemaVersion != 1 || e.StreamSequence < 1 || e.AggregateVersion < 1 || e.EventID == "" || e.OccurredAt.IsZero() || !identifier.MatchString(e.OrderID) || string(r.Key) != e.OrderID {
		return errors.New("invalid Orders event envelope")
	}
	switch e.Type {
	case "OrderCompleted", "OrderCancelled":
		return nil // terminal notifications are not payment requests
	case "OrderAwaitingPayment":
	default:
		return errors.New("unsupported Orders event")
	}
	var body struct {
		Dataset     string `json:"dataset_id"`
		Order       string `json:"order_id"`
		State       string `json:"state"`
		Version     int64  `json:"version"`
		Reservation struct {
			Amount   *int64 `json:"total_minor"`
			Currency string `json:"currency"`
		} `json:"reservation"`
	}
	if err := json.Unmarshal(e.Payload, &body); err != nil || body.Dataset != e.Dataset || body.Order != e.OrderID || body.Version != e.AggregateVersion || body.State != "PAYMENT_PENDING" || body.Reservation.Amount == nil || *body.Reservation.Amount < 0 || !currencyCode.MatchString(body.Reservation.Currency) {
		return errors.New("invalid reserved amount or order snapshot")
	}
	approved, delay := p.simulate(e.OrderID)
	return pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `INSERT INTO payments.jobs(order_id,event_id,envelope,amount_minor,currency,approve,ready_at)
 VALUES($1,$2,$3,$4,$5,$6,clock_timestamp()+$7*interval '1 millisecond') ON CONFLICT DO NOTHING`, e.OrderID, e.EventID, r.Value, *body.Reservation.Amount, body.Reservation.Currency, approved, delay.Milliseconds())
		if err != nil || tag.RowsAffected() == 1 {
			return err
		}
		var same bool
		if err = tx.QueryRow(ctx, "SELECT EXISTS(SELECT FROM payments.jobs WHERE order_id=$1 AND event_id=$2 AND envelope=$3::jsonb)", e.OrderID, e.EventID, r.Value).Scan(&same); err != nil {
			return err
		}
		if !same {
			return errors.New("conflicting payment request")
		}
		return nil
	})
}

// finish locks a due job only after the simulated delay. The debit, terminal
// state and transactional outbox are committed together or all roll back.
func (p processor) finish(ctx context.Context) error {
	return pgx.BeginFunc(ctx, p.pool, func(tx pgx.Tx) error {
		var order, cause, currency string
		var amount int64
		var approve bool
		err := tx.QueryRow(ctx, `SELECT order_id,event_id,amount_minor,currency,approve FROM payments.jobs
 WHERE state='PROCESSING' AND ready_at<=clock_timestamp() ORDER BY ready_at,order_id LIMIT 1 FOR UPDATE SKIP LOCKED`).Scan(&order, &cause, &amount, &currency, &approve)
		if err != nil {
			return err
		}
		state, kind := "DECLINED", "PaymentDeclined"
		if approve {
			if _, err = tx.Exec(ctx, "INSERT INTO payments.debits(order_id,amount_minor,currency) VALUES($1,$2,$3)", order, amount, currency); err != nil {
				return err
			}
			state, kind = "SUCCEEDED", "PaymentSucceeded"
		}
		if _, err = tx.Exec(ctx, "UPDATE payments.jobs SET state=$2 WHERE order_id=$1", order, state); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO payments.outbox(event_id,order_id,event_type,payload) VALUES($1,$2,$3,
 jsonb_build_object('order_id',$2::text,'state',$4::text,'amount_minor',$5::bigint,'currency',$6::text,'causation_id',$7::text))`, order+":payment-result", order, kind, state, amount, currency, cause)
		return err
	})
}
