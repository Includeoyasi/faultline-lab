-- The entire database belongs to ONE immutable dataset. A new dataset needs
-- a new isolated database/Compose project, never deletion of individual keys.
CREATE TABLE inventory.dataset (
    singleton boolean PRIMARY KEY DEFAULT true CHECK (singleton),
    dataset_id text NOT NULL CHECK (length(dataset_id) BETWEEN 1 AND 128)
);

CREATE TABLE inventory.reservations (
    order_id text PRIMARY KEY CHECK (length(order_id) BETWEEN 1 AND 128),
    state text NOT NULL CHECK (state IN ('RESERVED', 'COMMITTED', 'RELEASED')),
    version bigint NOT NULL CHECK (version > 0),
    total_minor bigint NOT NULL CHECK (total_minor >= 0),
    currency text NOT NULL CHECK (currency = '' OR currency ~ '^[A-Z]{3}$'),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE inventory.reservation_items (
    order_id text NOT NULL REFERENCES inventory.reservations(order_id),
    sku text NOT NULL REFERENCES inventory.products(sku),
    quantity bigint NOT NULL CHECK (quantity > 0),
    unit_price_minor bigint NOT NULL CHECK (unit_price_minor >= 0),
    PRIMARY KEY (order_id, sku)
);

CREATE TABLE inventory.operations (
    operation_id text PRIMARY KEY CHECK (length(operation_id) BETWEEN 1 AND 128),
    request_hash text NOT NULL CHECK (length(request_hash) = 64),
    receipt jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE inventory.ledger (
    order_id text NOT NULL REFERENCES inventory.reservations(order_id),
    version bigint NOT NULL,
    sku text NOT NULL REFERENCES inventory.products(sku),
    available_delta bigint NOT NULL,
    operation_id text NOT NULL REFERENCES inventory.operations(operation_id) DEFERRABLE INITIALLY DEFERRED,
    PRIMARY KEY (order_id, version, sku)
);

-- Durable local facts only. Publication and F09 fencing are a separate PR.
CREATE TABLE inventory.outbox (
    event_id text PRIMARY KEY,
    order_id text NOT NULL REFERENCES inventory.reservations(order_id),
    aggregate_version bigint NOT NULL CHECK (aggregate_version > 0),
    event_type text NOT NULL CHECK (event_type IN ('InventoryReserved', 'InventoryCommitted', 'InventoryReleased')),
    payload jsonb NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (order_id, aggregate_version)
);
