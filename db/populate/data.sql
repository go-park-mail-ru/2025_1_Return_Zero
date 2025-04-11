TRUNCATE TABLE track CASCADE;
TRUNCATE TABLE album CASCADE;
TRUNCATE TABLE artist CASCADE;
TRUNCATE TABLE genre CASCADE;
TRUNCATE TABLE genre_album CASCADE;
TRUNCATE TABLE genre_track CASCADE;
TRUNCATE TABLE track_artist CASCADE;

-- var tracks = []Track{
-- 	{ID: 1, Title: "Lagtrain", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-HdxXE6BxJ65FHooi-rtiaPw-t500x500.jpg", Album: "Anticyclone"},
-- 	{ID: 2, Title: "Lost Umbrella", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-Z9Jm9zLWMUzmOePX-TiOdqA-t500x500.jpg", Album: "Anticyclone"},
-- 	{ID: 3, Title: "Racing Into The Night", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-9fxbzFYK9QjT0aIg-eXpu8Q-t1080x1080.jpg", Album: "THE BOOK"},
-- 	{ID: 4, Title: "Idol", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-g677ppuycPRMga7w-LwVVlQ-t500x500.jpg", Album: "THE BOOK"},
-- 	{ID: 5, Title: "Monster", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-ztyGtBiqtACBb5zy-WtrLrg-t500x500.jpg", Album: "THE BOOK"},
-- 	{ID: 6, Title: "KICK BACK", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-lXWDlsG2J1UVytER-8YKCOg-t1080x1080.jpg", Album: "BOOTLEG"},
-- 	{ID: 7, Title: "Lemon", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-000446001171-xnyep8-t500x500.jpg", Album: "BOOTLEG"},
-- 	{ID: 8, Title: "Peace Sign", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-000482219301-jrnq0h-t500x500.jpg", Album: "BOOTLEG"},
-- 	{ID: 9, Title: "Sparkle", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000452912388-ft13zk-t1080x1080.jpg", Album: "Your Name."},
-- 	{ID: 10, Title: "Nandemonaiya", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000230768346-878y9o-t500x500.jpg", Album: "Your Name."},
-- 	{ID: 11, Title: "Suzume", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-OR55dgkv9l0JHg6J-NUMaSQ-t500x500.jpg", Album: "Your Name."},
-- 	{ID: 12, Title: "Pretender", Artist: "Official HIGE DANdism", Image: "https://i1.sndcdn.com/artworks-000644002372-j1fgr1-t500x500.jpg", Album: "Your Name."},
-- 	{ID: 13, Title: "Mixed Nuts", Artist: "Official HIGE DANdism", Image: "https://i1.sndcdn.com/artworks-68ZsJYMEYjCHMEpM-z4UHxg-t500x500.jpg", Album: "Your Name."},
-- 	{ID: 14, Title: "Cry Baby", Artist: "Official HIGE DANdism", Image: "https://i1.sndcdn.com/artworks-G0RPyB0xahP2CyHW-4H1THQ-t500x500.jpg", Album: "Your Name."},
-- 	{ID: 15, Title: "Dream Lantern", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000350712186-6xaoo7-t500x500.jpg", Album: "Your Name."},
-- 	{ID: 16, Title: "Zenzenzense", Artist: "RADWIMPS", Image: "https://i1.sndcdn.com/artworks-000189644938-tywci0-t1080x1080.jpg", Album: "Your Name."},
-- 	{ID: 17, Title: "Shinigami", Artist: "Kenshi Yonezu", Image: "https://i1.sndcdn.com/artworks-Z0nrZzzmeWrfD6ny-iVaI8w-t500x500.jpg", Album: "BOOTLEG"},
-- 	{ID: 18, Title: "Gunjo", Artist: "YOASOBI", Image: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSpTz3Ys6eJleSr2shfdB2BMq15WsipNe4rgQ&s", Album: "BOOTLEG"},
-- 	{ID: 19, Title: "Tabun", Artist: "YOASOBI", Image: "https://i1.sndcdn.com/artworks-dumxejUZ4jURPErm-xUFVFw-t500x500.jpg", Album: "BOOTLEG"},
-- 	{ID: 20, Title: "Ghost City Tokyo", Artist: "Inabakumori", Image: "https://i1.sndcdn.com/artworks-ssoxHlQypZXAQKap-tEfJ6A-t500x500.jpg", Album: "BOOTLEG"},
-- }

-- INSERT INTO artist (title, thumbnail_url, description) VALUES
-- 	('Inabakumori', 'https://i1.sndcdn.com/artworks-000640888066-bwv7e8-t500x500.jpg', 'Inabakumori is a Japanese artist'),
-- 	('YOASOBI', 'https://i.scdn.co/image/ab67616100005174bfdd8a29d0c6bc6950055234', 'YOASOBI is a Japanese artist'),
-- 	('Kenshi Yonezu', 'https://i.scdn.co/image/ab6761610000e5ebd7ca899f6e53b54976a8594b', 'Kenshi Yonezu is a Japanese artist'),
-- 	('RADWIMPS', 'https://i.scdn.co/image/ab6761610000e5ebc9d443fb5ced1dd32d106632', 'RADWIMPS is a Japanese artist'),
-- 	('Official HIGE DANdism', 'https://i.scdn.co/image/ab6761610000e5ebf9f7513528a90d1dde6d3aaa', 'Official HIGE DANdism is a Japanese artist');

-- INSERT INTO album (title, artist_id, thumbnail_url, release_date, type) VALUES
-- 	('Anticyclone', 1, 'https://i.scdn.co/image/ab67616d0000b27325c2a3af824b7dd8cafae97e', '2023-01-01', 'album'),
-- 	('THE BOOK', 2, 'https://i.scdn.co/image/ab67616d0000b273684d81c9356531f2a456b1c1', '2024-01-01', 'album'),
-- 	('BOOTLEG', 3, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQFG72O6ftYjIepEZw_aMvGYuE5kPvnll6v9g&s', '2022-01-01', 'ep'),
-- 	('Your Name.', 4, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQ0oNJ9dV6ldbzePBS8FsQcVoE3tPwEw3aqhw&s', '2021-01-01', 'album'),
--     ('Official HIGE DANdism', 5, 'https://i.scdn.co/image/ab6761610000e5ebf9f7513528a90d1dde6d3aaa', '2020-01-01', 'album');

INSERT INTO genre (name) VALUES
	('J-Pop'),
	('Rock'),
	('Electronic'),
	('Pop'),
	('Hip-Hop');

-- INSERT INTO genre_album (genre_id, album_id) VALUES
-- 	(1, 1),
-- 	(2, 1),
-- 	(3, 2),
-- 	(4, 3),
-- 	(5, 4);

-- INSERT INTO track (title, album_id, duration, thumbnail_url, file_url, position) VALUES
--     ('Lagtrain', 1, 216, 'https://i1.sndcdn.com/artworks-HdxXE6BxJ65FHooi-rtiaPw-t500x500.jpg', '', 1),
--     ('Lost Umbrella', 1, 216, 'https://i1.sndcdn.com/artworks-Z9Jm9zLWMUzmOePX-TiOdqA-t500x500.jpg', '', 2),
--     ('Racing Into The Night', 2, 216, 'https://i1.sndcdn.com/artworks-9fxbzFYK9QjT0aIg-eXpu8Q-t1080x1080.jpg', '', 1),
--     ('Idol', 2, 216, 'https://i1.sndcdn.com/artworks-g677ppuycPRMga7w-LwVVlQ-t500x500.jpg', '', 2),
--     ('Monster', 2, 216, 'https://i1.sndcdn.com/artworks-ztyGtBiqtACBb5zy-WtrLrg-t500x500.jpg', 'https://ia801503.us.archive.org/12/items/Lagtrain/Lagtrain.mp3', 3),
--     ('KICK BACK', 3, 216, 'https://i1.sndcdn.com/artworks-lXWDlsG2J1UVytER-8YKCOg-t1080x1080.jpg', '', 1),
--     ('Lemon', 3, 216, 'https://i1.sndcdn.com/artworks-000446001171-xnyep8-t500x500.jpg', '', 2),
--     ('Peace Sign', 3, 216, 'https://i1.sndcdn.com/artworks-000482219301-jrnq0h-t500x500.jpg', '', 3),
--     ('Sparkle', 4, 216, 'https://i1.sndcdn.com/artworks-000452912388-ft13zk-t1080x1080.jpg', '', 1),
--     ('Nandemonaiya', 4, 216, 'https://i1.sndcdn.com/artworks-000230768346-878y9o-t500x500.jpg', '', 2),
--     ('Suzume', 4, 216, 'https://i1.sndcdn.com/artworks-OR55dgkv9l0JHg6J-NUMaSQ-t500x500.jpg', '', 3),
--     ('Pretender', 5, 216, 'https://i1.sndcdn.com/artworks-000644002372-j1fgr1-t500x500.jpg', '', 1),
--     ('Mixed Nuts', 5, 216, 'https://i1.sndcdn.com/artworks-68ZsJYMEYjCHMEpM-z4UHxg-t500x500.jpg', '', 2),
--     ('Cry Baby', 5, 216, 'https://i1.sndcdn.com/artworks-G0RPyB0xahP2CyHW-4H1THQ-t500x500.jpg', '', 3),
--     ('Dream Lantern', 5, 216, 'https://i1.sndcdn.com/artworks-000350712186-6xaoo7-t500x500.jpg', '', 4),
--     ('Zenzenzense', 5, 216, 'https://i1.sndcdn.com/artworks-000189644938-tywci0-t1080x1080.jpg', '', 5),
--     ('Shinigami', 5, 216, 'https://i1.sndcdn.com/artworks-Z0nrZzzmeWrfD6ny-iVaI8w-t500x500.jpg', '', 6),
--     ('Gunjo', 5, 216, 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcSpTz3Ys6eJleSr2shfdB2BMq15WsipNe4rgQ&s', '', 7),
--     ('Tabun', 5, 216, 'https://i1.sndcdn.com/artworks-dumxejUZ4jURPErm-xUFVFw-t500x500.jpg', '', 8),
--     ('Ghost City Tokyo', 5, 216, 'https://i1.sndcdn.com/artworks-ssoxHlQypZXAQKap-tEfJ6A-t500x500.jpg', '', 9);


-- INSERT INTO genre_track (genre_id, track_id) VALUES
--     (1, 1),
--     (2, 1),
--     (3, 2),
--     (4, 3),
--     (5, 4);

-- INSERT INTO track_artist (track_id, artist_id, role) VALUES
--     (1, 1, 'main'),
--     (2, 1, 'main'),
--     (3, 2, 'main'),
--     (4, 2, 'main'),
--     (5, 2, 'main'),
--     (6, 3, 'main'),
--     (7, 3, 'main'),
--     (8, 3, 'main'),
--     (9, 4, 'main'),
--     (10, 4, 'main'),
--     (11, 4, 'main'),
--     (12, 5, 'main'),
--     (13, 5, 'main'),
--     (14, 5, 'main'),
--     (15, 5, 'main'),
--     (16, 5, 'main'),
--     (17, 5, 'main'),
--     (18, 5, 'main'),
--     (19, 5, 'main'),
--     (20, 5, 'main');
