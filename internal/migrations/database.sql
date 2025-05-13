
-- USERS
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- ACCOUNTS
CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    number VARCHAR(34) NOT NULL,
    balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    currency VARCHAR(8) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_account_number ON accounts(number);

-- CARDS
CREATE TABLE IF NOT EXISTS cards (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    number VARCHAR(19) NOT NULL,
    expiry_month INT NOT NULL CHECK (expiry_month >= 1 AND expiry_month <= 12),
    expiry_year INT NOT NULL,
    encrypted_cvv VARCHAR(255) NOT NULL,
    cardholder_name VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX IF NOT EXISTS uidx_card_number ON cards(number);

-- TRANSACTIONS
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    account_id BIGINT NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
    amount NUMERIC(18,2) NOT NULL,
    currency VARCHAR(8) NOT NULL,
    type VARCHAR(32) NOT NULL,     -- deposit, withdraw, transfer, payment, etc.
    status VARCHAR(32) NOT NULL,   -- success, pending, failed, etc.
    description TEXT,
    related_entity_id BIGINT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- CREDITS
CREATE TABLE IF NOT EXISTS credits (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    amount NUMERIC(18,2) NOT NULL,
    currency VARCHAR(8) NOT NULL,
    rate NUMERIC(8,3) NOT NULL,           -- процентная ставка
    term_months INT NOT NULL,
    status VARCHAR(32) NOT NULL,          -- active, closed, overdue и т.д.
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- PAYMENT_SCHEDULES
CREATE TABLE IF NOT EXISTS payment_schedules (
    id BIGSERIAL PRIMARY KEY,
    credit_id BIGINT NOT NULL REFERENCES credits(id) ON DELETE CASCADE,
    payment_num INT NOT NULL,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    is_paid BOOLEAN NOT NULL DEFAULT FALSE,
    paid_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_payment_schedule_credit_id ON payment_schedules(credit_id);
