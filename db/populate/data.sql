TRUNCATE TABLE track CASCADE;
TRUNCATE TABLE album CASCADE;
TRUNCATE TABLE artist CASCADE;
TRUNCATE TABLE genre CASCADE;
TRUNCATE TABLE genre_album CASCADE;
TRUNCATE TABLE genre_track CASCADE;
TRUNCATE TABLE track_artist CASCADE;

INSERT INTO artist (title, thumbnail_url, description) VALUES
	('Inabakumori', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/anticyclone.jpg', 'Inabakumori is a Japanese artist'),
	('YOASOBI', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/yoasobi.jpg', 'YOASOBI is a Japanese artist'),
	('Kenshi Yonezu', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/kenshiyonezu.jpg', 'Kenshi Yonezu is a Japanese artist'),
	('RADWIMPS', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/radwimps.jpg', 'RADWIMPS is a Japanese artist'),
	('Official HIGE DANdism', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/officialhigedandism.jpg', 'Official HIGE DANdism is a Japanese artist'),
    ('Toaka', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/toaka.jpg', 'Toaka is a Japanese artist'),
    ('Ayase', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/ayase.jpg', 'Ayase is a Japanese artist');

INSERT INTO album (title, thumbnail_url, release_date, type) VALUES
	('Anticyclone', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/anticyclone.jpg', '2023-01-01', 'album'),
	('THE BOOK', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/thebook.jpg', '2024-01-01', 'album'),
	('BOOTLEG', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/bootleg.jpg', '2022-01-01', 'ep'),
	('Your Name.', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/yourname.jpg', '2021-01-01', 'album'),
    ('Official HIGE DANdism', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/officialhigedandism.jpg', '2020-01-01', 'album'),
    ('Ghost City Tokyo', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/ghostcitytokyo.jpg', '2024-01-01', 'album');

INSERT INTO genre (name) VALUES
	('J-Pop'),
	('Rock'),
	('Electronic'),
	('Pop'),
	('Hip-Hop');

INSERT INTO genre_album (genre_id, album_id) VALUES
	(1, 1),
	(2, 1),
	(3, 2),
	(4, 3),
	(5, 4);

INSERT INTO track (title, album_id, duration, thumbnail_url, file_url, position) VALUES
    ('Lagtrain', 1, 252, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/lagtrain.jpg', 'lagtrain.mp3', 1),
    ('Lost Umbrella', 1, 255, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/lostumbrella.jpg', 'lostumbrella.mp3', 2),
    ('Racing Into The Night', 2, 275, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/racingintotheright.jpg', 'racingintotheright.mp3', 1),
    ('Idol', 2, 226, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/idol.jpg', 'idol.mp3', 2),
    ('Monster', 2, 206, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/monster.jpg', 'monster.mp3', 3),
    ('KICK BACK', 3, 194, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/kickback.jpg', 'kickback.mp3', 1),
    ('Lemon', 3, 256, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/lemon.jpg', 'lemon.mp3', 2),
    ('Peace Sign', 3, 237, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/peacesign.jpg', 'peacesign.mp3', 3),
    ('Sparkle', 4, 534, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/sparkle.jpg', 'sparkle.mp3', 1),
    ('Nandemonaiya', 4, 344, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/nandemonaiya.jpg', 'nandemonaiya.mp3', 2),
    ('Suzume', 4, 236, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/suzume.jpg', 'suzume.mp3', 3),
    ('Pretender', 5, 327, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/pretender.jpg', 'pretender.mp3', 1),
    ('Mixed Nuts', 5, 213, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/mixednuts.jpg', 'mixednuts.mp3', 2),
    ('Cry Baby', 5, 242, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/crybaby.jpg', 'crybaby.mp3', 3),
    ('Dream Lantern', 4, 129, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/dreamlantern.jpg', 'dreamlantern.mp3', 4),
    ('Zenzenzense', 4, 277, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/zenzenzense.jpg', 'zenzenzense.mp3', 5),
    ('Shinigami', 3, 181, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/shinigami.jpg', 'shinigami.mp3', 4),
    ('Gunjo', 2, 248, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/gunjo.jpg', 'gunjo.mp3', 4),
    ('Tabun', 2, 263, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/tabun.jpg', 'tabun.mp3', 5),
    ('Ghost City Tokyo', 6, 204, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/ghostcitytokyo.jpg', 'ghostcitytokyo.mp3', 1);

INSERT INTO genre_track (genre_id, track_id) VALUES
    (1, 1),
    (2, 1),
    (3, 2),
    (4, 3),
    (5, 4);

INSERT INTO track_artist (track_id, artist_id, role) VALUES
    (1, 1, 'main'),
    (2, 1, 'main'),
    (3, 2, 'main'),
    (4, 2, 'main'),
    (5, 2, 'main'),
    (6, 3, 'main'),
    (7, 3, 'main'),
    (8, 3, 'main'),
    (9, 4, 'main'),
    (10, 4, 'main'),
    (11, 4, 'main'),
    (11, 6, 'featured'),
    (12, 5, 'main'),
    (13, 5, 'main'),
    (14, 5, 'main'),
    (15, 4, 'main'),
    (16, 4, 'main'),
    (17, 3, 'main'),
    (18, 2, 'main'),
    (19, 2, 'main'),
    (20, 7, 'main');

INSERT INTO album_artist (album_id, artist_id) VALUES
    (1, 1),
    (2, 2),
    (3, 3),
    (4, 4),
    (5, 5),
    (6, 7);

REFRESH MATERIALIZED VIEW artist_stats;
REFRESH MATERIALIZED VIEW album_stats;
REFRESH MATERIALIZED VIEW track_stats;
