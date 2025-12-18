CREATE DATABASE assignment;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    email VARCHAR(100) UNIQUE
);

CREATE TABLE rewards (
    id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id),
    stock_symbol VARCHAR(20) NOT NULL,
    quantity NUMERIC(18, 6) NOT NULL,
    rewarded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ledger (
    id SERIAL PRIMARY KEY,
    reward_id INT REFERENCES rewards(id),
    account_type VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    amount NUMERIC(18, 4) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stock_prices (
    id SERIAL PRIMARY KEY,
    stock_symbol VARCHAR(20),
    price_inr NUMERIC(18, 4),
    captured_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (name, email) VALUES ('John Doe', 'john@example.com');