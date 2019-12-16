-- +migrate Up

CREATE TABLE store.modules (
	id             TEXT      NOT NULL,
	user_category  TEXT      NOT NULL,
	type           TEXT      NOT NULL,
	version        BIGINT    NOT NULL,
	data           JSONB     NOT NULL,
	PRIMARY KEY (id, user_category)
);

CREATE TABLE store.modules_history (
	version        BIGSERIAL,
	id             TEXT      NOT NULL,
	user_category  TEXT      NOT NULL,
	type           TEXT      NOT NULL,
	data           JSONB,
	deleted        BOOLEAN   NOT NULL DEFAULT false,
	created_at     TIMESTAMP DEFAULT current_timestamp,
	PRIMARY KEY (version)
);

-- +migrate Down

