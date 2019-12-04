-- +migrate Up

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

-- +migrate Down

DROP TABLE users.user_param_log;