-- +migrate Up

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