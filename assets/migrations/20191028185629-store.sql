-- +migrate Up

CREATE TABLE store (
        id SERIAL,
        name VARCHAR(250),
        status INT                                       NOT NULL DEFAULT 1,
        tenant_id INT                                    NOT NULL,
        created_at TIMESTAMPTZ DEFAULT current_timestamp,
        updated_at TIMESTAMPTZ DEFAULT current_timestamp,
        PRIMARY KEY (id)
    );

-- +migrate Down

DROP TABLE store;