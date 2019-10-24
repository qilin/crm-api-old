DO $$DECLARE game_id UUID;
BEGIN
	INSERT INTO games (
		tenant_id,
		published,
		name,
		description,
		-- tagline,
		release_date,
		copyright,
		genre
	)
	VALUES (
		1, -- owner organization
		true, 
		'Northgard',
		'Northgard is a strategy game based on Norse mythology in which you control a clan of Vikings vying for the control of a mysterious newfound continent.',
		'2018.03.07',
		'(C) 2016 Shiro Games. The Shiro Games name and logo and the Northgard name and logo are trademarks of Shiro Games and may be registered trademarks in certain countries. All rights reserved.',
		'Indie'
	) 
	RETURNING id INTO game_id;

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
		(game_id, name)
	VALUES
		(game_id, 'Shiro Games');

	INSERT INTO game_publishers
		(game_id, name)
	VALUES
		(game_id, 'Shiro Games');

	INSERT INTO game_platforms
		(game_id, platform)
	VALUES
		(game_id, 'linux'),
		(game_id, 'win'),
		(game_id, 'mac');

	INSERT INTO game_links
		(game_id, site)
	VALUES
		(game_id, 'http://northgard.net/');

	INSERT INTO game_languages
		(game_id, language, interface, sound, subtitles)
	VALUES
		(game_id, 'English',             true,  true, true),
		(game_id, 'French',              true, false, true),
		(game_id, 'German',              true, false, true),
		(game_id, 'Russian',             true, false, true),
		(game_id, 'Simplified Chinese',  true, false, true),
		(game_id, 'Portuguese - Brazil', true, false, true),
		(game_id, 'Polish',              true, false, true),
		(game_id, 'Spanish - Spain',     true, false, true),
		(game_id, 'Turkish',             true, false, true);


	INSERT INTO game_thumbnails
		(game_id, format, url)
	VALUES
		(game_id, '2x1', 'https://steamcdn-a.akamaihd.net/steam/apps/466560/header.jpg?t=1571828499');

	INSERT INTO game_images
		(game_id, width, height, url)
	VALUES
		(game_id, 1920, 1080, 'https://steamcdn-a.akamaihd.net/steam/apps/466560/ss_0de6fea3166a04027b671e438645e98186095c63.jpg?t=1571828499');

END$$;