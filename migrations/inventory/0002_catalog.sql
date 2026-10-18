CREATE TABLE inventory.products (
    sku text PRIMARY KEY CHECK (length(sku) BETWEEN 1 AND 64),
    title text NOT NULL CHECK (length(title) BETWEEN 1 AND 200),
    price_minor bigint NOT NULL CHECK (price_minor >= 0),
    currency text NOT NULL CHECK (currency ~ '^[A-Z]{3}$')
);

CREATE TABLE inventory.stock (
    sku text PRIMARY KEY REFERENCES inventory.products(sku),
    available bigint NOT NULL CHECK (available >= 0)
);
