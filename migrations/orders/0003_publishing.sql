-- Allocate numbers under the same row lock and transaction as the business fact.
LOCK TABLE orders.outbox IN ACCESS EXCLUSIVE MODE;
ALTER TABLE orders.outbox ADD COLUMN stream_sequence bigint,
    ADD COLUMN published_at timestamptz;
WITH numbered AS (
    SELECT event_id, row_number() OVER (PARTITION BY order_id ORDER BY aggregate_version) AS seq
    FROM orders.outbox
)
UPDATE orders.outbox o SET stream_sequence = n.seq
FROM numbered n WHERE n.event_id = o.event_id;
ALTER TABLE orders.outbox ALTER COLUMN stream_sequence SET NOT NULL,
    ADD CHECK (stream_sequence > 0), ADD UNIQUE (order_id, stream_sequence);

CREATE TABLE orders.outbox_counter (
    order_id text PRIMARY KEY,
    last_sequence bigint NOT NULL CHECK (last_sequence > 0)
);
INSERT INTO orders.outbox_counter
SELECT order_id, max(stream_sequence) FROM orders.outbox GROUP BY order_id;

CREATE FUNCTION orders.number_outbox_event() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
    INSERT INTO orders.outbox_counter AS c VALUES (NEW.order_id, 1)
    ON CONFLICT (order_id) DO UPDATE SET last_sequence = c.last_sequence + 1
    RETURNING last_sequence INTO NEW.stream_sequence;
    RETURN NEW;
END
$$;
CREATE TRIGGER number_outbox_event BEFORE INSERT ON orders.outbox
FOR EACH ROW EXECUTE FUNCTION orders.number_outbox_event();
CREATE INDEX unpublished_stream ON orders.outbox (order_id, stream_sequence)
WHERE published_at IS NULL;

CREATE TABLE orders.publisher_lease (
    singleton boolean PRIMARY KEY CHECK (singleton),
    generation bigint NOT NULL CHECK (generation >= 0),
    lease_until timestamptz NOT NULL
);
INSERT INTO orders.publisher_lease VALUES (true, 0, '-infinity');
