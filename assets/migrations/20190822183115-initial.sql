-- +migrate Up

CREATE TABLE list
(
    id         BIGSERIAL                              NOT NULL,
    name       TEXT,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp,
    PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE list;