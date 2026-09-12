CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS accounts (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    owner_id UUID NOT NULL,
    currency VARCHAR(3) NOT NULL,
    type VARCHAR NOT NULL DEFAULT 'AVAILABLE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT unique_owner_currency UNIQUE (owner_id, currency)
);

CREATE INDEX IF NOT EXISTS index_accounts_owner
ON accounts
USING HASH (owner_id);
