-- +migrate Up

CREATE TABLE actions_log (
		id         BIGSERIAL NOT NULL,
		data	   JSONB     NOT NULL,
        created_at TIMESTAMP DEFAULT current_timestamp,
        updated_at TIMESTAMP DEFAULT current_timestamp,
        PRIMARY KEY (id)
    );

-- +migrate Down

DROP TABLE actions_log;