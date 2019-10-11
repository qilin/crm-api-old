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

INSERT INTO jwt_keys (alg, iss, key_type, key) VALUES('ES256','P1','pem', '-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEEVs/o5+uQbTjL3chynL4wXgUg2R9
q9UU8I5mEovUf86QZ7kOBIjJwqnzD1omageEHWwHdBO6B+dFabmdT9POxg==
-----END PUBLIC KEY-----');

-- +migrate StatementEnd

-- +migrate Down

DROP TABLE jwt_keys;