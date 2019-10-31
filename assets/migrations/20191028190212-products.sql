-- +migrate Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE products (
        id          uuid DEFAULT uuid_generate_v4 (),
        platform_id bigint                                NOT NULL,
        url         TEXT                                   NOT NULL,
        created_at  TIMESTAMPTZ DEFAULT current_timestamp,
        updated_at  TIMESTAMPTZ DEFAULT current_timestamp,
        PRIMARY KEY (id)
    );

-- +migrate Down
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP TABLE products;