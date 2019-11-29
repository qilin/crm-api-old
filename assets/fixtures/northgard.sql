-- +migrate Up

-- +migrate StatementBegin
DO $$DECLARE game_id UUID;
BEGIN
	INSERT INTO games (
		tenant_id,
		published,
		release_date,
		copyright,
		genre,
		default_lang
	)
	VALUES (
		1, -- owner organization
		true, 
		'2018.03.07',
		'(C) 2016 Shiro Games. The Shiro Games name and logo and the Northgard name and logo are trademarks of Shiro Games and may be registered trademarks in certain countries. All rights reserved.',
		'Indie',
		'en'
	) 
	RETURNING id INTO game_id;

	INSERT INTO game_props (
		id,
		lang,
		name,
		description,
		-- tagline,
		site
	)
	VALUES (
		game_id,
		'en',
		'Northgard',
		'Northgard is a strategy game based on Norse mythology in which you control a clan of Vikings vying for the control of a mysterious newfound continent.',
		'http://northgard.net/'
	);

	INSERT INTO game_genres 
		(game_id, genre)
	VALUES 
		(game_id, 'Strategy'),
		(game_id, 'Simulation');

	INSERT INTO game_features
		(game_id, feature)
	VALUES
		(game_id, 'Singleplayer'),
		(game_id, 'Multiplayer'),
		(game_id, 'Achivements');

	INSERT INTO game_developers
		(game_id, name, lang)
	VALUES
		(game_id, 'Shiro Games', 'en');

	INSERT INTO game_publishers
		(game_id, name, lang)
	VALUES
		(game_id, 'Shiro Games', 'en');

	INSERT INTO game_platforms
		(game_id, platform)
	VALUES
		(game_id, 'linux'),
		(game_id, 'win'),
		(game_id, 'mac');

	INSERT INTO game_links
		(game_id, lang, name, url)
	VALUES
		(game_id, 'en', 'facebook', 'https://www.facebook.com/northgard/');

	INSERT INTO game_languages
		(game_id, language, interface, sound, subtitles)
	VALUES
		(game_id, 'en',         true,  true, true),
		(game_id, 'fr',         true, false, true),
		(game_id, 'de',         true, false, true),
		(game_id, 'ru',         true, false, true),
		(game_id, 'zh',         true, false, true), -- TODO add zh-Hans
		(game_id, 'pt-BR',      true, false, true),
		(game_id, 'pl',         true, false, true),
		(game_id, 'es-ES',      true, false, true),
		(game_id, 'tr',         true, false, true);


	INSERT INTO game_thumbnails
		(game_id, format, url)
	VALUES
		(game_id, '2x1', 'https://steamcdn-a.akamaihd.net/steam/apps/466560/header.jpg?t=1571828499');

	INSERT INTO game_images
		(game_id, lang, width, height, url)
	VALUES
		(game_id, 'en', 1920, 1080, 'https://steamcdn-a.akamaihd.net/steam/apps/466560/ss_0de6fea3166a04027b671e438645e98186095c63.jpg?t=1571828499');

	INSERT INTO products
		(codename, type, object_id, tenant_id)
	VALUES
		('northgard', 'Game', game_id, 1);

END$$;
-- +migrate StatementEnd

-- +migrate Down
