CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "accounts" (
    "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
    "owner_id" UUID NOT NULL,
    "currency" varchar(3) NOT NULL,
    type VARCHAR NOT NULL DEFAULT 'AVAILABLE',
    "created_at" timestamptz NOT NULL DEFAULT NOW()
    "updated_at" timestamptz NOT NULL DEFAULT NOW()

    CONSTRAINT unique_owner_currency UNIQUE (owner_id, currency);
);

CREATE INDEX IF NOT EXISTS index_accounts_owner ON "accounts" (owner_id) 
USING HASH (owner_id);