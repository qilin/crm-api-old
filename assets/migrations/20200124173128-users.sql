-- +migrate Up

CREATE TABLE IF NOT EXISTS public.users
(
    id         serial      NOT NULL,
    tenant_id  INT         NOT NULL,
    status     boolean     NOT NULL,
    email      text        NOT NULL,
    picture    text        NOT NULL,
    first_name varchar     NOT NULL,
    last_name  varchar     NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (id)
);

-- +migrate Down

DROP TABLE public.users;