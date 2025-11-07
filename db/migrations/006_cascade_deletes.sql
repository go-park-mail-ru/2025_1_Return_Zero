-- Add cascading foreign keys to ensure dependent rows are cleaned up automatically

-- Remove orphaned rows so new constraints can be validated
DELETE FROM track_artist ta WHERE NOT EXISTS (
    SELECT 1 FROM track t WHERE t.id = ta.track_id
);

DELETE FROM album_artist aa WHERE NOT EXISTS (
    SELECT 1 FROM album a WHERE a.id = aa.album_id
);

DELETE FROM playlist_track pt WHERE NOT EXISTS (
    SELECT 1 FROM track t WHERE t.id = pt.track_id
);

DELETE FROM genre_track gt WHERE NOT EXISTS (
    SELECT 1 FROM track t WHERE t.id = gt.track_id
);

DELETE FROM genre_album ga WHERE NOT EXISTS (
    SELECT 1 FROM album a WHERE a.id = ga.album_id
);

DELETE FROM track t WHERE NOT EXISTS (
    SELECT 1 FROM album a WHERE a.id = t.album_id
);

ALTER TABLE track
ADD CONSTRAINT fk_track_album
FOREIGN KEY (album_id) REFERENCES album (id)
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE track_artist
ADD CONSTRAINT fk_track_artist_track
FOREIGN KEY (track_id) REFERENCES track (id)
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE album_artist
ADD CONSTRAINT fk_album_artist_album
FOREIGN KEY (album_id) REFERENCES album (id)
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE playlist_track
ADD CONSTRAINT fk_playlist_track_track
FOREIGN KEY (track_id) REFERENCES track (id)
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE genre_track
ADD CONSTRAINT fk_genre_track_track
FOREIGN KEY (track_id) REFERENCES track (id)
ON DELETE CASCADE
ON UPDATE CASCADE;

ALTER TABLE genre_album
ADD CONSTRAINT fk_genre_album_album
FOREIGN KEY (album_id) REFERENCES album (id)
ON DELETE CASCADE
ON UPDATE CASCADE;

---- create above / drop below ----

ALTER TABLE genre_album
DROP CONSTRAINT IF EXISTS fk_genre_album_album;

ALTER TABLE genre_track
DROP CONSTRAINT IF EXISTS fk_genre_track_track;

ALTER TABLE playlist_track
DROP CONSTRAINT IF EXISTS fk_playlist_track_track;

ALTER TABLE album_artist
DROP CONSTRAINT IF EXISTS fk_album_artist_album;

ALTER TABLE track_artist
DROP CONSTRAINT IF EXISTS fk_track_artist_track;

ALTER TABLE track
DROP CONSTRAINT IF EXISTS fk_track_album;
