-- +migrate Up

CREATE UNIQUE INDEX ux_games_slug ON store.games((data->>'slug'));

-- +migrate Down

DROP INDEX store.ux_games_slug;

