INSERT INTO inventory.products (sku, title, price_minor, currency) VALUES
    ('SKU-001', 'Notebook', 19900, 'RUB'),
    ('SKU-002', 'Coffee', 45000, 'RUB'),
    ('SKU-003', 'Headphones', 399000, 'RUB')
ON CONFLICT (sku) DO NOTHING;

INSERT INTO inventory.stock (sku, available) VALUES
    ('SKU-001', 100), ('SKU-002', 50), ('SKU-003', 25)
ON CONFLICT (sku) DO NOTHING;
