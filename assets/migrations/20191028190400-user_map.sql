-- +migrate Up

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE user_map (
        user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        store_id INT NOT NULL,
        external_id TEXT NOT NULL,
        created_at TIMESTAMPTZ DEFAULT current_timestamp,
        updated_at TIMESTAMPTZ DEFAULT current_timestamp
    );

-- +migrate Down

DROP TABLE user_map;

DROP EXTENSION pgcrypto;