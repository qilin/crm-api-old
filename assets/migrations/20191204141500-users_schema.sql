-- +migrate Up

CREATE SCHEMA IF NOT EXISTS users;

-- +migrate Down

DROP SCHEMA users;