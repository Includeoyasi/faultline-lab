-- Each service database belongs to one dataset. Monetary values use minor units.
CREATE TABLE payments.debits (
    order_id text PRIMARY KEY CHECK (length(order_id) BETWEEN 1 AND 128),
    amount_minor bigint NOT NULL CHECK (amount_minor >= 0),
    currency text NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    recorded_at timestamptz NOT NULL DEFAULT clock_timestamp()
);

-- One full refund references the immutable debit. Its amount and currency
-- are derived from that debit, so callers cannot supply a different price.
CREATE TABLE payments.refunds (
    order_id text PRIMARY KEY REFERENCES payments.debits(order_id),
    recorded_at timestamptz NOT NULL DEFAULT clock_timestamp()
);

CREATE VIEW payments.payment_ledger AS
    SELECT order_id, 'DEBIT'::text AS entry_type, amount_minor, currency, recorded_at
    FROM payments.debits
    UNION ALL
    SELECT r.order_id, 'REFUND'::text, d.amount_minor, d.currency, r.recorded_at
    FROM payments.refunds r JOIN payments.debits d USING (order_id);

-- Append-only ledger: corrections are separate entries.
CREATE FUNCTION payments.reject_ledger_change() RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    RAISE EXCEPTION 'payment ledger entries are immutable' USING ERRCODE = '23514';
END;
$$;
CREATE TRIGGER immutable_debits BEFORE UPDATE OR DELETE ON payments.debits
FOR EACH ROW EXECUTE FUNCTION payments.reject_ledger_change();
CREATE TRIGGER immutable_refunds BEFORE UPDATE OR DELETE ON payments.refunds
FOR EACH ROW EXECUTE FUNCTION payments.reject_ledger_change();
