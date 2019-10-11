-- +migrate Up

CREATE TABLE jwt_keys
(
    id         BIGSERIAL                              NOT NULL,
    alg        VARCHAR(6)                             NOT NULL,
    iss        TEXT                                   NOT NULL,
    key_type   VARCHAR(3)                             NOT NULL,
    key        TEXT                                   NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT current_timestamp,
    updated_at TIMESTAMPTZ DEFAULT current_timestamp,
    PRIMARY KEY (id)
);

-- +migrate StatementBegin

INSERT INTO jwt_keys (alg, iss, key_type, key) VALUES('ES256','P1','pem', '-----BEGIN PUBLIC KEY-----\n      MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEa9fBxVOSbv7iQ3UC0xMwzdqK3QGU\n      1uGDpgNhuOGdNVdleb5iYfGyYqvXPWN02gwFBePLWYBKEPslgeUQEpJ0GQ==\n      -----END PUBLIC KEY-----');

-- +migrate StatementEnd

-- +migrate Down

DROP TABLE jwt_keys;