-- +migrate Up

CREATE TABLE users
(
    id         BIGSERIAL                              NOT NULL,
    name       varchar(255)                           NOT NULL UNIQUE,
    password   varchar(70),
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp,
    PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE list;