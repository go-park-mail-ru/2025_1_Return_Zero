-- Write your migrate up statements here

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
	((SELECT id FROM genre WHERE name = 'J-Pop'), (SELECT id FROM album WHERE title = 'Anticyclone')),
	((SELECT id FROM genre WHERE name = 'Rock'), (SELECT id FROM album WHERE title = 'Anticyclone')),
	((SELECT id FROM genre WHERE name = 'Electronic'), (SELECT id FROM album WHERE title = 'THE BOOK')),
	((SELECT id FROM genre WHERE name = 'Pop'), (SELECT id FROM album WHERE title = 'BOOTLEG')),
	((SELECT id FROM genre WHERE name = 'Hip-Hop'), (SELECT id FROM album WHERE title = 'Your Name.'));

INSERT INTO track (title, album_id, duration, thumbnail_url, file_url, position) VALUES
    ('Lagtrain', (SELECT id FROM album WHERE title = 'Anticyclone'), 252, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/lagtrain.jpg', 'lagtrain.mp3', 1),
    ('Lost Umbrella', (SELECT id FROM album WHERE title = 'Anticyclone'), 255, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/lostumbrella.jpg', 'lostumbrella.mp3', 2),
    ('Racing Into The Night', (SELECT id FROM album WHERE title = 'THE BOOK'), 275, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/racingintotheright.jpg', 'racingintotheright.mp3', 1),
    ('Idol', (SELECT id FROM album WHERE title = 'THE BOOK'), 226, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/idol.jpg', 'idol.mp3', 2),
    ('Monster', (SELECT id FROM album WHERE title = 'THE BOOK'), 206, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/monster.jpg', 'monster.mp3', 3),
    ('KICK BACK', (SELECT id FROM album WHERE title = 'BOOTLEG'), 194, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/kickback.jpg', 'kickback.mp3', 1),
    ('Lemon', (SELECT id FROM album WHERE title = 'BOOTLEG'), 256, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/lemon.jpg', 'lemon.mp3', 2),
    ('Peace Sign', (SELECT id FROM album WHERE title = 'BOOTLEG'), 237, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/peacesign.jpg', 'peacesign.mp3', 3),
    ('Sparkle', (SELECT id FROM album WHERE title = 'Your Name.'), 534, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/sparkle.jpg', 'sparkle.mp3', 1),
    ('Nandemonaiya', (SELECT id FROM album WHERE title = 'Your Name.'), 344, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/nandemonaiya.jpg', 'nandemonaiya.mp3', 2),
    ('Suzume', (SELECT id FROM album WHERE title = 'Your Name.'), 236, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/suzume.jpg', 'suzume.mp3', 3),
    ('Pretender', (SELECT id FROM album WHERE title = 'Official HIGE DANdism'), 327, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/pretender.jpg', 'pretender.mp3', 1),
    ('Mixed Nuts', (SELECT id FROM album WHERE title = 'Official HIGE DANdism'), 213, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/mixednuts.jpg', 'mixednuts.mp3', 2),
    ('Cry Baby', (SELECT id FROM album WHERE title = 'Official HIGE DANdism'), 242, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/crybaby.jpg', 'crybaby.mp3', 3),
    ('Dream Lantern', (SELECT id FROM album WHERE title = 'Your Name.'), 129, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/dreamlantern.jpg', 'dreamlantern.mp3', 4),
    ('Zenzenzense', (SELECT id FROM album WHERE title = 'Your Name.'), 277, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/zenzenzense.jpg', 'zenzenzense.mp3', 5),
    ('Shinigami', (SELECT id FROM album WHERE title = 'BOOTLEG'), 181, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/shinigami.jpg', 'shinigami.mp3', 4),
    ('Gunjo', (SELECT id FROM album WHERE title = 'THE BOOK'), 248, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/gunjo.jpg', 'gunjo.mp3', 4),
    ('Tabun', (SELECT id FROM album WHERE title = 'THE BOOK'), 263, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/tabun.jpg', 'tabun.mp3', 5),
    ('Ghost City Tokyo', (SELECT id FROM album WHERE title = 'Ghost City Tokyo'), 204, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/ghostcitytokyo.jpg', 'ghostcitytokyo.mp3', 1);


INSERT INTO genre_track (genre_id, track_id) VALUES
    ((SELECT id FROM genre WHERE name = 'J-Pop'), (SELECT id FROM track WHERE title = 'Lagtrain')),
    ((SELECT id FROM genre WHERE name = 'Rock'), (SELECT id FROM track WHERE title = 'Lagtrain')),
    ((SELECT id FROM genre WHERE name = 'Electronic'), (SELECT id FROM track WHERE title = 'Lost Umbrella')),
    ((SELECT id FROM genre WHERE name = 'Pop'), (SELECT id FROM track WHERE title = 'Racing Into The Night')),
    ((SELECT id FROM genre WHERE name = 'Hip-Hop'), (SELECT id FROM track WHERE title = 'Idol'));

INSERT INTO track_artist (track_id, artist_id, role) VALUES
    ((SELECT id FROM track WHERE title = 'Lagtrain'), (SELECT id FROM artist WHERE title = 'Inabakumori'), 'main'),
    ((SELECT id FROM track WHERE title = 'Lost Umbrella'), (SELECT id FROM artist WHERE title = 'Inabakumori'), 'main'),
    ((SELECT id FROM track WHERE title = 'Racing Into The Night'), (SELECT id FROM artist WHERE title = 'YOASOBI'), 'main'),
    ((SELECT id FROM track WHERE title = 'Idol'), (SELECT id FROM artist WHERE title = 'YOASOBI'), 'main'),
    ((SELECT id FROM track WHERE title = 'Monster'), (SELECT id FROM artist WHERE title = 'YOASOBI'), 'main'),
    ((SELECT id FROM track WHERE title = 'KICK BACK'), (SELECT id FROM artist WHERE title = 'Kenshi Yonezu'), 'main'),
    ((SELECT id FROM track WHERE title = 'Lemon'), (SELECT id FROM artist WHERE title = 'Kenshi Yonezu'), 'main'),
    ((SELECT id FROM track WHERE title = 'Peace Sign'), (SELECT id FROM artist WHERE title = 'Kenshi Yonezu'), 'main'),
    ((SELECT id FROM track WHERE title = 'Sparkle'), (SELECT id FROM artist WHERE title = 'RADWIMPS'), 'main'),
    ((SELECT id FROM track WHERE title = 'Nandemonaiya'), (SELECT id FROM artist WHERE title = 'RADWIMPS'), 'main'),
    ((SELECT id FROM track WHERE title = 'Suzume'), (SELECT id FROM artist WHERE title = 'RADWIMPS'), 'main'),
    ((SELECT id FROM track WHERE title = 'Suzume'), (SELECT id FROM artist WHERE title = 'Toaka'), 'featured'),
    ((SELECT id FROM track WHERE title = 'Pretender'), (SELECT id FROM artist WHERE title = 'Official HIGE DANdism'), 'main'),
    ((SELECT id FROM track WHERE title = 'Mixed Nuts'), (SELECT id FROM artist WHERE title = 'Official HIGE DANdism'), 'main'),
    ((SELECT id FROM track WHERE title = 'Cry Baby'), (SELECT id FROM artist WHERE title = 'Official HIGE DANdism'), 'main'),
    ((SELECT id FROM track WHERE title = 'Dream Lantern'), (SELECT id FROM artist WHERE title = 'RADWIMPS'), 'main'),
    ((SELECT id FROM track WHERE title = 'Zenzenzense'), (SELECT id FROM artist WHERE title = 'RADWIMPS'), 'main'),
    ((SELECT id FROM track WHERE title = 'Shinigami'), (SELECT id FROM artist WHERE title = 'Kenshi Yonezu'), 'main'),
    ((SELECT id FROM track WHERE title = 'Gunjo'), (SELECT id FROM artist WHERE title = 'YOASOBI'), 'main'),
    ((SELECT id FROM track WHERE title = 'Tabun'), (SELECT id FROM artist WHERE title = 'YOASOBI'), 'main'),
    ((SELECT id FROM track WHERE title = 'Ghost City Tokyo'), (SELECT id FROM artist WHERE title = 'Ayase'), 'main');

INSERT INTO album_artist (album_id, artist_id) VALUES
    ((SELECT id FROM album WHERE title = 'Anticyclone'), (SELECT id FROM artist WHERE title = 'Inabakumori')),
    ((SELECT id FROM album WHERE title = 'THE BOOK'), (SELECT id FROM artist WHERE title = 'YOASOBI')),
    ((SELECT id FROM album WHERE title = 'BOOTLEG'), (SELECT id FROM artist WHERE title = 'Kenshi Yonezu')),
    ((SELECT id FROM album WHERE title = 'Your Name.'), (SELECT id FROM artist WHERE title = 'RADWIMPS')),
    ((SELECT id FROM album WHERE title = 'Official HIGE DANdism'), (SELECT id FROM artist WHERE title = 'Official HIGE DANdism')),
    ((SELECT id FROM album WHERE title = 'Ghost City Tokyo'), (SELECT id FROM artist WHERE title = 'Ayase'));

REFRESH MATERIALIZED VIEW artist_stats;
REFRESH MATERIALIZED VIEW album_stats;
REFRESH MATERIALIZED VIEW track_stats;

---- create above / drop below ----

TRUNCATE TABLE track CASCADE;
TRUNCATE TABLE album CASCADE;
TRUNCATE TABLE artist CASCADE;
TRUNCATE TABLE genre CASCADE;
TRUNCATE TABLE genre_album CASCADE;
TRUNCATE TABLE genre_track CASCADE;
TRUNCATE TABLE track_artist CASCADE;
TRUNCATE TABLE album_artist CASCADE;
TRUNCATE TABLE album_genre CASCADE;
TRUNCATE TABLE track_genre CASCADE;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
