package database

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Exercise the actual application grants. All fixture rows roll back together.
func checkPaymentLedger(t *testing.T, ctx context.Context, conn *pgx.Conn) {
	t.Helper()
	tx, err := conn.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer rollback(tx)
	run := func(sql string) {
		t.Helper()
		if _, err := tx.Exec(ctx, sql); err != nil {
			t.Fatal(err)
		}
	}
	reject := func(sql, code string) {
		t.Helper()
		nested, err := tx.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		_, err = nested.Exec(ctx, sql)
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgErr.Code != code {
			t.Fatalf("expected SQLSTATE %s, got %v", code, err)
		}
		if err = nested.Rollback(ctx); err != nil {
			t.Fatal(err)
		}
	}
	reject("INSERT INTO payments.refunds(order_id) VALUES ('ledger-test')", "23503")
	run("INSERT INTO payments.debits(order_id,amount_minor,currency) VALUES ('ledger-test',1250,'RUB')")
	reject("INSERT INTO payments.debits(order_id,amount_minor,currency) VALUES ('ledger-test',1300,'RUB')", "23505")
	reject("UPDATE payments.debits SET amount_minor=1 WHERE order_id='ledger-test'", "23514")
	run("INSERT INTO payments.refunds(order_id) VALUES ('ledger-test')")
	reject("INSERT INTO payments.refunds(order_id) VALUES ('ledger-test')", "23505")
	reject("DELETE FROM payments.refunds WHERE order_id='ledger-test'", "23514")
	var debits, refunds int
	err = tx.QueryRow(ctx, `SELECT count(*) FILTER (WHERE entry_type='DEBIT'),count(*) FILTER (WHERE entry_type='REFUND')
 FROM payments.payment_ledger WHERE order_id='ledger-test' AND amount_minor=1250 AND currency='RUB'`).Scan(&debits, &refunds)
	if err != nil || debits != 1 || refunds != 1 {
		t.Fatalf("debits=%d refunds=%d error=%v", debits, refunds, err)
	}
}
