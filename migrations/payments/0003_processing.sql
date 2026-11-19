CREATE TABLE payments.dataset (
    singleton boolean PRIMARY KEY DEFAULT true CHECK (singleton),
    dataset_id text NOT NULL
);

-- Transactional inbox and durable job share a row. Keeping the original
-- envelope lets us distinguish replay from a reused identity.
CREATE TABLE payments.jobs (
    order_id text PRIMARY KEY,
    event_id text NOT NULL UNIQUE,
    envelope jsonb NOT NULL,
    amount_minor bigint NOT NULL CHECK (amount_minor >= 0),
    currency text NOT NULL CHECK (currency ~ '^[A-Z]{3}$'),
    approve boolean NOT NULL,
    ready_at timestamptz NOT NULL,
    state text NOT NULL DEFAULT 'PROCESSING'
        CHECK (state IN ('PROCESSING','SUCCEEDED','DECLINED'))
);
CREATE INDEX ready_payments ON payments.jobs (ready_at,order_id) WHERE state='PROCESSING';

-- Each charge has one terminal result in this milestone.
CREATE TABLE payments.outbox (
    event_id text PRIMARY KEY,
    order_id text UNIQUE NOT NULL REFERENCES payments.jobs,
    aggregate_version bigint NOT NULL DEFAULT 2 CHECK (aggregate_version=2),
    stream_sequence bigint NOT NULL DEFAULT 1 CHECK (stream_sequence=1),
    event_type text NOT NULL CHECK (event_type IN ('PaymentSucceeded','PaymentDeclined')),
    payload jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT clock_timestamp(),
    published_at timestamptz
);
CREATE INDEX pending_results ON payments.outbox(order_id) WHERE published_at IS NULL;
CREATE TABLE payments.publisher_lease (
    singleton boolean PRIMARY KEY CHECK (singleton),
    generation bigint NOT NULL CHECK (generation >= 0),
    lease_until timestamptz NOT NULL
);
INSERT INTO payments.publisher_lease VALUES(true,0,'-infinity');
