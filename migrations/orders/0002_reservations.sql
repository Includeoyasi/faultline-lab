-- One isolated database is bound to one dataset; protective records have no TTL.
CREATE TABLE orders.dataset (
    singleton boolean PRIMARY KEY DEFAULT true CHECK (singleton),
    dataset_id text NOT NULL CHECK (length(dataset_id) BETWEEN 1 AND 128)
);
CREATE TABLE orders.orders (
    order_id text PRIMARY KEY,
    request_key text NOT NULL UNIQUE,
    request_hash text NOT NULL CHECK (length(request_hash) = 64),
    reserve_operation_id text NOT NULL UNIQUE,
    state text NOT NULL CHECK (state IN ('ACCEPTED','RESERVING','PAYMENT_PENDING','REJECTED')),
    version bigint NOT NULL CHECK (version > 0),
    items jsonb NOT NULL CHECK (jsonb_typeof(items) = 'array' AND jsonb_array_length(items) BETWEEN 1 AND 100),
    reservation jsonb,
    rejection text NOT NULL DEFAULT '',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CHECK ((state = 'PAYMENT_PENDING') = (reservation IS NOT NULL)),
    CHECK ((state = 'REJECTED') = (rejection <> ''))
);
CREATE TABLE orders.reserve_jobs (
    order_id text PRIMARY KEY REFERENCES orders.orders(order_id),
    generation bigint NOT NULL DEFAULT 0 CHECK (generation >= 0),
    order_version bigint NOT NULL DEFAULT 1 CHECK (order_version > 0),
    attempts bigint NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    lease_until timestamptz,
    available_at timestamptz NOT NULL DEFAULT now(),
    done boolean NOT NULL DEFAULT false,
    blocked_reason text NOT NULL DEFAULT '',
    last_error text NOT NULL DEFAULT ''
);
CREATE INDEX reserve_jobs_ready ON orders.reserve_jobs(available_at,order_id)
    WHERE NOT done AND blocked_reason = '';
-- Local facts only. Kafka publication and F09 are a separate package.
CREATE TABLE orders.outbox (
    event_id text PRIMARY KEY,
    order_id text NOT NULL REFERENCES orders.orders(order_id),
    aggregate_version bigint NOT NULL CHECK (aggregate_version > 0),
    event_type text NOT NULL CHECK (event_type = 'OrderAwaitingPayment'),
    payload jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (order_id,aggregate_version)
);
