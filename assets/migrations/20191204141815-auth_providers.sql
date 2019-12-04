-- +migrate Up

CREATE TABLE users.authentication_providers (
-- user_id
    user_id BIGINT NOT NULL REFERENCES users.users(id),
-- provider
    provider VARCHAR(32) NOT NULL ,
    provider_id text NOT NULL,
    provider_key text NOT NULL,
-- ts
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp,

-- keys & indexes
    PRIMARY KEY (user_id, provider, provider_id)
);

-- +migrate Down

DROP TABLE users.authentication_providers;