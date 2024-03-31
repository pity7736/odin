BEGIN;

CREATE TABLE users(
    id UUID PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE tokens(
    id SERIAL PRIMARY KEY,
    value TEXT NOT NULL UNIQUE,
    user_id UUID REFERENCES users
);


CREATE TABLE IF NOT EXISTS categories(
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(15) NOT NULL,
    user_id UUID REFERENCES users
);

CREATE TABLE IF NOT EXISTS wallets(
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    balance NUMERIC(12, 2) NOT NULL,
    user_id UUID REFERENCES users NOT NULL
);

CREATE TYPE movement_types AS ENUM ('I', 'E');

CREATE TABLE IF NOT EXISTS movements(
    id UUID PRIMARY KEY,
    amount NUMERIC(10, 2) NOT NULL,
    date date NOT NULL,
    movement_type movement_types NOT NULL,
    wallet_id UUID REFERENCES wallets NOT NULL,
    category_id UUID REFERENCES categories NOT NULL
);

CREATE TABLE IF NOT EXISTS transfers(
    id UUID PRIMARY KEY,
    amount NUMERIC(10, 2) NOT NULL,
    date date NOT NULL,
    source_id UUID REFERENCES wallets NOT NULL,
    target_id UUID REFERENCES wallets NOT NULL,
    expense_id UUID REFERENCES movements NOT NULL,
    income_id UUID REFERENCES movements NOT NULL
);

COMMIT;
