--
-- SNIPPETS
--

INSERT INTO snippets (title, content, created, expires) VALUES (
	'An old silent pond',
	'An old silent pond...\nA frog jump into the pond,\nsplash! Silence again.\n\n- Matsuo Basho',
	UTC_TIMESTAMP(),
	DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created, expires) VALUES (
	'Over the wintry forest',
	'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n- Natsume Soseki',
	UTC_TIMESTAMP(),
	DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title, content, created, expires) VALUES (
	'First autumn morning',
	'First autumn morning\nthe mirror I stare into\nshows my fathers face.\n\n- Murakami Kijo',
	UTC_TIMESTAMP(),
	DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

--
-- USERS
--

INSERT INTO users(name, email, hashed_password, created) VALUES (
	'Alice Jones',
	'alice@example.com',
	'$2a$12$3gxV7F216Z7prgzEAqwTuu8930hoPPp/md.EQn6E4UUoCdz/VEaYe',
	'2022-01-01 10:00:00'
);
