CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE role AS ENUM(
    'member',
    'admin'
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    user_role role NOT NULL, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE products (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    discount NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (
        discount >= 0
        AND discount <= price
    ),
    description TEXT,
    image TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE receipts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    sum_price NUMERIC(12, 2) NOT NULL CHECK (sum_price >= 0),
    sum_discount NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (
        sum_discount >= 0
        AND sum_discount <= sum_price
    ),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE receipt_products (
    receipt_id BIGINT NOT NULL REFERENCES receipts (id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products (id) ON DELETE RESTRICT,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    price NUMERIC(12, 2) NOT NULL CHECK (price >= 0),
    discount NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (
        discount >= 0
        AND discount <= price
    ),
    PRIMARY KEY (receipt_id, product_id)
);

CREATE INDEX idx_receipts_user_id ON receipts (user_id);

CREATE INDEX idx_receipt_products_product_id ON receipt_products (product_id);