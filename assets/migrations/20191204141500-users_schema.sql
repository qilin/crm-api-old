-- +migrate Up

CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE users.users (
-- id
    id          BIGSERIAL NOT NULL,
-- login & password
    email       VARCHAR(256),
    phone       VARCHAR(31),
    password    VARCHAR(65) NOT NULL DEFAULT '',
-- account status &
    status          INT8, -- active, blocked, locked, deleted, etc.
    service_level   INT8, -- service grade: bronze, gold, prior, vip, etc.
-- address
    address_1   VARCHAR,
    address_2   VARCHAR,
    city        VARCHAR,
    state       VARCHAR,
    country     VARCHAR,
    zip         VARCHAR(20),
-- name & dob
    photo_url   TEXT,
    first_name  VARCHAR(100),
    last_name   VARCHAR(100),
    birth_date  INT,
-- localisation
    language    VARCHAR,
-- ts
    created_at TIMESTAMP DEFAULT current_timestamp,
    updated_at TIMESTAMP DEFAULT current_timestamp,
-- keys & indexes
    PRIMARY KEY (id)
);

CREATE TABLE users.user_providers_map (
-- user_id
    user_id         BIGINT      NOT NULL REFERENCES users.users(id),
-- provider
    provider        VARCHAR(32) NOT NULL ,
    provider_id     TEXT        NOT NULL,
    provider_key    TEXT        NOT NULL,
-- ts
    created_at      TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at      TIMESTAMPTZ DEFAULT current_timestamp,
-- keys & indexes
    PRIMARY KEY (user_id, provider, provider_id)
);

CREATE TABLE users.user_param_log (
-- id
    id         BIGSERIAL NOT NULL,
-- fk
    user_id    BIGINT    NOT NULL REFERENCES users.users(id),
-- log
    param	   VARCHAR   NOT NULL,
    old        VARCHAR,
    new        VARCHAR,
-- log
    updated_by VARCHAR, -- who updated the param: user, support, etc.
    user_agent VARCHAR,
    ip         VARCHAR,
    hw_id       VARCHAR, -- hardware id
-- ts
    created_at TIMESTAMP DEFAULT current_timestamp,
-- keys & indexes
    PRIMARY KEY (id)
);

CREATE TABLE users.auth_log (
-- id
    id          BIGSERIAL NOT NULL,
-- user
    user_id     BIGINT  NOT NULL REFERENCES users.users(id),
-- log
    action      VARCHAR NOT NULL, -- login / logout / sign up
    user_agent  VARCHAR,
    ip          VARCHAR,
    hw_id       VARCHAR, -- hardware id
-- ts
    created_at TIMESTAMP DEFAULT current_timestamp,
-- keys & indexes
    PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE users.auth_log;
DROP TABLE users.user_param_log;
DROP TABLE users.user_providers_map;
DROP TABLE users.users;
DROP SCHEMA users;