-- Add cascading foreign keys to ensure dependent rows are cleaned up automatically

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
