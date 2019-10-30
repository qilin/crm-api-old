-- +migrate Up

CREATE TABLE platforms (
        id SERIAL,
        name VARCHAR(250),
        status INT                                       NOT NULL DEFAULT 1,
        created_at TIMESTAMPTZ DEFAULT current_timestamp,
        updated_at TIMESTAMPTZ DEFAULT current_timestamp,
        PRIMARY KEY (id)
    );

-- +migrate Down

DROP TABLE platforms;