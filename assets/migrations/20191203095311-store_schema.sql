-- +migrate Up

CREATE SCHEMA store;

CREATE TABLE store.games (
	id         UUID      NOT NULL,
	version    BIGINT    NOT NULL,
	data       JSONB     NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE store.games_history (
	id         UUID      NOT NULL,
	version    BIGSERIAL,
	data       JSONB,
	deleted    BOOLEAN   NOT NULL DEFAULT false,
	created_at TIMESTAMP DEFAULT current_timestamp,
	PRIMARY KEY (version)
);

-- +migrate Down

DROP TABLE store.games_history;
DROP TABLE store.games;
DROP SCHEMA store;