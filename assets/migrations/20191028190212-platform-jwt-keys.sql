-- +migrate Up

CREATE TABLE platform_jwt_keys (
        id         BIGSERIAL                              NOT NULL,
        alg        VARCHAR(6)                             NOT NULL,
        iss        TEXT                                   NOT NULL,
        key_type   VARCHAR(3)                             NOT NULL,
        key        TEXT                                   NOT NULL UNIQUE,
        created_at TIMESTAMPTZ DEFAULT current_timestamp,
        updated_at TIMESTAMPTZ DEFAULT current_timestamp,
        PRIMARY KEY (id)
    );

-- +migrate Down

DROP TABLE platform_jwt_keys;