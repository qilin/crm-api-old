-- +migrate Up

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

-- +migrate Down

DROP TABLE users.users;