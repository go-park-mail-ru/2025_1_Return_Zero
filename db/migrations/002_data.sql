-- Write your migrate up statements here

INSERT INTO artist (title, thumbnail_url, description) VALUES
	('Inabakumori', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/anticyclone.jpg', 'Inabakumori is a Japanese artist'),
	('YOASOBI', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/yoasobi.jpg', 'YOASOBI is a Japanese artist'),
	('Kenshi Yonezu', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/kenshiyonezu.jpg', 'Kenshi Yonezu is a Japanese artist'),
	('RADWIMPS', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/radwimps.jpg', 'RADWIMPS is a Japanese artist'),
	('Official HIGE DANdism', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/officialhigedandism.jpg', 'Official HIGE DANdism is a Japanese artist'),
    ('Toaka', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/toaka.jpg', 'Toaka is a Japanese artist'),
    ('Ayase', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/ayase.jpg', 'Ayase is a Japanese artist'),
    ('Eminem', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/eminem.jpg', 'Cool man'),
    ('The Cab', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/thecab.jpg', 'Angel'),
    ('Sergey Eybog', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/sergeyeybog.jpg', 'Free musician'),
    ('Katalepsy', 'https://returnzeroimages.fra1.digitaloceanspaces.com/artists/katalepsy.jpg', 'Slamming brutal death metal band from Russia');

INSERT INTO album (title, thumbnail_url, release_date, type) VALUES
	('Anticyclone', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/anticyclone.jpg', '2023-01-01', 'album'),
	('THE BOOK', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/thebook.jpg', '2024-01-01', 'album'),
	('BOOTLEG', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/bootleg.jpg', '2022-01-01', 'ep'),
	('Your Name.', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/yourname.jpg', '2021-01-01', 'album'),
    ('Official HIGE DANdism', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/officialhigedandism.jpg', '2020-01-01', 'album'),
    ('Ghost City Tokyo', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/ghostcitytokyo.jpg', '2024-01-01', 'album'),
    ('The Eminem Show', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/theeminemshow.jpg', '2002-01-01', 'album'),
    ('Symphony Soldier', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/symphonysoldier.jpg', '2015-01-01', 'single'),
    ('Everlasting Summer', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/everlastingsummer.jpg', '2016-01-01', 'album'),
    ('Music Brings Injures', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/musicbringsinjures.webp', '2024-01-01', 'album'),
    ('Triumph Of Evilution', 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/triumphofevilution.webp', '2024-01-01', 'ep');

INSERT INTO genre (name) VALUES
	('J-Pop'),
	('Rock'),
	('Electronic'),
	('Pop'),
	('Hip-Hop'),
    ('OST');

INSERT INTO genre_album (genre_id, album_id) VALUES
	((SELECT id FROM genre WHERE name = 'J-Pop'), (SELECT id FROM album WHERE title = 'Anticyclone')),
	((SELECT id FROM genre WHERE name = 'Electronic'), (SELECT id FROM album WHERE title = 'THE BOOK')),
	((SELECT id FROM genre WHERE name = 'Pop'), (SELECT id FROM album WHERE title = 'BOOTLEG')),
	((SELECT id FROM genre WHERE name = 'OST'), (SELECT id FROM album WHERE title = 'Your Name.')),
    ((SELECT id FROM genre WHERE name = 'Hip-Hop'), (SELECT id FROM album WHERE title = 'The Eminem Show')),
    ((SELECT id FROM genre WHERE name = 'Rock'), (SELECT id FROM album WHERE title = 'Symphony Soldier')),
    ((SELECT id FROM genre WHERE name = 'OST'), (SELECT id FROM album WHERE title = 'Everlasting Summer'));

INSERT INTO track (title, album_id, duration, thumbnail_url, file_url, position) VALUES
    ('Gialo', (SELECT id FROM album WHERE title = 'Music Brings Injures'), 212, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/musicbringsinjures.webp', 'gialo.mp3', 1),
    ('Sluggish Cranial Grinding', (SELECT id FROM album WHERE title = 'Music Brings Injures'), 224, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/musicbringsinjures.webp', 'sluggishcranialgrinding.mp3', 2),
    ('Rabid', (SELECT id FROM album WHERE title = 'Music Brings Injures'), 137, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/musicbringsinjures.webp', 'rabid.mp3', 3),
    ('Necroviolated To Liquid', (SELECT id FROM album WHERE title = 'Music Brings Injures'), 226, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/musicbringsinjures.webp', 'necroviolatedtoliquid.mp3', 4),
    ('ConsumingTheAbyss', (SELECT id FROM album WHERE title = 'Music Brings Injures'), 223, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/musicbringsinjures.webp', 'consumingtheabyss.mp3', 5),
    ('S.O.D.', (SELECT id FROM album WHERE title = 'Music Brings Injures'), 265, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/musicbringsinjures.webp', 'sod.mp3', 6),
    ('Post-Apocalyptic Segregation', (SELECT id FROM album WHERE title = 'Triumph Of Evilution'), 251, 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/triumphofevilution.webp', 'postapocalypticsegregation.mp3', 1),
    ('Carpet Wounding', (SELECT id FROM album WHERE title = 'Triumph Of Evilution'), 255, 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/triumphofevilution.webp', 'carpetwounding.mp3', 2),
    ('H. Tearing Sinew', (SELECT id FROM album WHERE title = 'Triumph Of Evilution'), 285, 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/triumphofevilution.webp', 'htearingsinew.mp3', 3),
    ('Number Of Death', (SELECT id FROM album WHERE title = 'Triumph Of Evilution'), 244, 'https://returnzeroimages.fra1.digitaloceanspaces.com/albums/triumphofevilution.webp', 'numberofdeath.mp3', 4),
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
    ('Ghost City Tokyo', (SELECT id FROM album WHERE title = 'Ghost City Tokyo'), 204, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/ghostcitytokyo.jpg', 'ghostcitytokyo.mp3', 1),
    ('Without Me', (SELECT id FROM album WHERE title = 'The Eminem Show'), 289, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/theeminemshow.jpg', 'withoutme.mp3', 1),
    ('Sing For The Moment', (SELECT id FROM album WHERE title = 'The Eminem Show'), 340, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/theeminemshow.jpg', 'singforthemoment.mp3', 2),
    ('Angel With A Shotgun', (SELECT id FROM album WHERE title = 'Symphony Soldier'), 201, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/symphonysoldier.jpg', 'angelwithashotgun.mp3', 1),
    ('Timid Girl', (SELECT id FROM album WHERE title = 'Everlasting Summer'), 82, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/everlastingsummer.jpg', 'timidgirl.mp3', 1),
    ('Let`s be friends', (SELECT id FROM album WHERE title = 'Everlasting Summer'), 123, 'https://returnzeroimages.fra1.digitaloceanspaces.com/tracks/everlastingsummer.jpg', 'letbefriends.mp3', 2);
    


INSERT INTO genre_track (genre_id, track_id) VALUES
    ((SELECT id FROM genre WHERE name = 'J-Pop'), (SELECT id FROM track WHERE title = 'Lagtrain')),
    ((SELECT id FROM genre WHERE name = 'Rock'), (SELECT id FROM track WHERE title = 'Lagtrain')),
    ((SELECT id FROM genre WHERE name = 'Electronic'), (SELECT id FROM track WHERE title = 'Lost Umbrella')),
    ((SELECT id FROM genre WHERE name = 'Pop'), (SELECT id FROM track WHERE title = 'Racing Into The Night')),
    ((SELECT id FROM genre WHERE name = 'Hip-Hop'), (SELECT id FROM track WHERE title = 'Idol')),
    ((SELECT id FROM genre WHERE name = 'OST'), (SELECT id FROM track WHERE title = 'Timid Girl')),
    ((SELECT id FROM genre WHERE name = 'OST'), (SELECT id FROM track WHERE title = 'Let`s be friends')),
    ((SELECT id FROM genre WHERE name = 'Rock'), (SELECT id FROM track WHERE title = 'Angel With A Shotgun'));

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
    ((SELECT id FROM track WHERE title = 'Ghost City Tokyo'), (SELECT id FROM artist WHERE title = 'Ayase'), 'main'),
    ((SELECT id FROM track WHERE title = 'Without Me'), (SELECT id FROM artist WHERE title = 'Eminem'), 'main'),
    ((SELECT id FROM track WHERE title = 'Sing For The Moment'), (SELECT id FROM artist WHERE title = 'Eminem'), 'main'),
    ((SELECT id FROM track WHERE title = 'Angel With A Shotgun'), (SELECT id FROM artist WHERE title = 'The Cab'), 'main'),
    ((SELECT id FROM track WHERE title = 'Timid Girl'), (SELECT id FROM artist WHERE title = 'Sergey Eybog'), 'main'),
    ((SELECT id FROM track WHERE title = 'Let`s be friends'), (SELECT id FROM artist WHERE title = 'Sergey Eybog'), 'main'),
    ((SELECT id FROM track WHERE title = 'Gialo'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main'),
    ((SELECT id FROM track WHERE title = 'Sluggish Cranial Grinding'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main'),
    ((SELECT id FROM track WHERE title = 'Rabid'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main'),
    ((SELECT id FROM track WHERE title = 'Necroviolated To Liquid'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main'),
    ((SELECT id FROM track WHERE title = 'ConsumingTheAbyss'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main'),
    ((SELECT id FROM track WHERE title = 'S.O.D.'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main'),
    ((SELECT id FROM track WHERE title = 'Post-Apocalyptic Segregation'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main'),
    ((SELECT id FROM track WHERE title = 'Carpet Wounding'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main'),
    ((SELECT id FROM track WHERE title = 'H. Tearing Sinew'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main'),
    ((SELECT id FROM track WHERE title = 'Number Of Death'), (SELECT id FROM artist WHERE title = 'Katalepsy'), 'main');


INSERT INTO album_artist (album_id, artist_id) VALUES
    ((SELECT id FROM album WHERE title = 'Anticyclone'), (SELECT id FROM artist WHERE title = 'Inabakumori')),
    ((SELECT id FROM album WHERE title = 'THE BOOK'), (SELECT id FROM artist WHERE title = 'YOASOBI')),
    ((SELECT id FROM album WHERE title = 'BOOTLEG'), (SELECT id FROM artist WHERE title = 'Kenshi Yonezu')),
    ((SELECT id FROM album WHERE title = 'Your Name.'), (SELECT id FROM artist WHERE title = 'RADWIMPS')),
    ((SELECT id FROM album WHERE title = 'Official HIGE DANdism'), (SELECT id FROM artist WHERE title = 'Official HIGE DANdism')),
    ((SELECT id FROM album WHERE title = 'Ghost City Tokyo'), (SELECT id FROM artist WHERE title = 'Ayase')),
    ((SELECT id FROM album WHERE title = 'The Eminem Show'), (SELECT id FROM artist WHERE title = 'Eminem')),
    ((SELECT id FROM album WHERE title = 'Symphony Soldier'), (SELECT id FROM artist WHERE title = 'The Cab')),
    ((SELECT id FROM album WHERE title = 'Everlasting Summer'), (SELECT id FROM artist WHERE title = 'Sergey Eybog')),
    ((SELECT id FROM album WHERE title = 'Music Brings Injures'), (SELECT id FROM artist WHERE title = 'Katalepsy')),
    ((SELECT id FROM album WHERE title = 'Triumph Of Evilution'), (SELECT id FROM artist WHERE title = 'Katalepsy'));

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

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
