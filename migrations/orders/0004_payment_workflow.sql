-- Existing PAYMENT_PENDING rows retain their reservation and can resume.
ALTER TABLE orders.orders
    DROP CONSTRAINT orders_state_check,
    DROP CONSTRAINT orders_check,
    ADD CONSTRAINT orders_state_check CHECK (state IN
        ('ACCEPTED','RESERVING','PAYMENT_PENDING','FINALIZING','COMPLETED','CANCELLING','CANCELLED','REJECTED')),
    ADD CONSTRAINT orders_reservation_check CHECK
        ((state IN ('PAYMENT_PENDING','FINALIZING','COMPLETED','CANCELLING','CANCELLED')) = (reservation IS NOT NULL));

ALTER TABLE orders.outbox DROP CONSTRAINT outbox_event_type_check,
    ADD CONSTRAINT outbox_event_type_check CHECK
        (event_type IN ('OrderAwaitingPayment','OrderCompleted','OrderCancelled'));

-- Transactional inbox. The cursor records the complete contiguous prefix
-- for each producer and order independently of business application.
CREATE TABLE orders.inbox_cursor (
    owner text NOT NULL CHECK (owner IN ('payments','inventory')),
    order_id text NOT NULL,
    sequence bigint NOT NULL CHECK (sequence >= 0),
    PRIMARY KEY (owner,order_id)
);
CREATE TABLE orders.inbox (
    owner text NOT NULL CHECK (owner IN ('payments','inventory')),
    event_id text NOT NULL,
    order_id text NOT NULL,
    stream_sequence bigint NOT NULL CHECK (stream_sequence > 0),
    body bytea NOT NULL,
    applied boolean NOT NULL DEFAULT false,
    available_at timestamptz NOT NULL DEFAULT clock_timestamp(),
    last_error text NOT NULL DEFAULT '',
    PRIMARY KEY (owner,event_id),
    UNIQUE (owner,order_id,stream_sequence)
);
CREATE INDEX inbox_pending ON orders.inbox(available_at,owner,order_id,stream_sequence)
    WHERE NOT applied;

-- The stable operation ID is reused across timeouts, crashes and replicas.
CREATE TABLE orders.inventory_actions (
    order_id text PRIMARY KEY REFERENCES orders.orders(order_id),
    operation_id text NOT NULL UNIQUE,
    action text NOT NULL CHECK (action IN ('commit','release')),
    generation bigint NOT NULL DEFAULT 0 CHECK (generation >= 0),
    attempts bigint NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    lease_until timestamptz,
    available_at timestamptz NOT NULL DEFAULT clock_timestamp(),
    done boolean NOT NULL DEFAULT false,
    blocked_reason text NOT NULL DEFAULT '',
    last_error text NOT NULL DEFAULT ''
);
CREATE INDEX inventory_actions_ready ON orders.inventory_actions(available_at,order_id)
    WHERE NOT done AND blocked_reason='';
