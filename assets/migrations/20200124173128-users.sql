-- +migrate Up

CREATE TABLE IF NOT EXISTS public.users
(
    id         serial      NOT NULL,
    tenant_id  INT         DEFAULT 1,
    status     boolean     DEFAULT true,
    email      text        NOT NULL,
    picture    text        DEFAULT '',
    first_name varchar     DEFAULT '',
    last_name  varchar     DEFAULT '',
    role       varchar     DEFAULT '',
    external_id varchar DEFAULT '',
    password    varchar DEFAULT '',
    auth_timestamp timestamptz NOT NULL DEFAULT now(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE public.users;