-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS pg_cron;

-- user microservice
CREATE TABLE IF NOT EXISTS "user" (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    CONSTRAINT email_length_check CHECK (LENGTH(email) >= 5 AND LENGTH(email) <= 30),
    username TEXT NOT NULL UNIQUE,
    CONSTRAINT username_length_check CHECK (LENGTH(username) >= 3 AND LENGTH(username) <= 20),
    thumbnail_url TEXT NOT NULL DEFAULT '/default_avatar.png',
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- user microservice
CREATE TABLE IF NOT EXISTS user_settings (
    user_id BIGINT NOT NULL PRIMARY KEY,
    is_public_playlists BOOLEAN NOT NULL DEFAULT FALSE, 
    is_public_minutes_listened BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_favorite_artists BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_tracks_listened BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_favorite_tracks BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_artists_listened BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
        REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- genre microservice
CREATE TABLE IF NOT EXISTS genre (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT genre_name_length_check CHECK (LENGTH(name) >= 3 AND LENGTH(name) <= 20)
);

-- artist microservice
CREATE TABLE IF NOT EXISTS artist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    CONSTRAINT artist_title_length_check CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
    description TEXT NOT NULL DEFAULT '',
    CONSTRAINT artist_description_length_check CHECK (LENGTH(description) <= 1000),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    thumbnail_url TEXT NOT NULL DEFAULT '/default_artist.png'
);

CREATE TABLE IF NOT EXISTS album (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    CONSTRAINT album_title_length_check CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
    type TEXT NOT NULL DEFAULT 'single',
    thumbnail_url TEXT NOT NULL DEFAULT '/default_album.png',
    release_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT album_valid_type_check CHECK (type IN ('album', 'single', 'ep', 'compilation'))
);

-- track microservice
CREATE TABLE IF NOT EXISTS track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    CONSTRAINT track_title_length_check CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
    thumbnail_url TEXT NOT NULL DEFAULT '/default_track.png',
    file_url TEXT NOT NULL DEFAULT '',
    album_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    duration INTEGER NOT NULL,
    position INTEGER NOT NULL,
    CONSTRAINT track_valid_duration_check CHECK (duration > 0),
    CONSTRAINT unique_album_track_check UNIQUE (album_id, position)
);

-- artist microservice
CREATE TABLE IF NOT EXISTS track_artist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, 
    track_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    role TEXT NOT NULL DEFAULT 'main',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (artist_id) 
        REFERENCES artist (id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    CONSTRAINT track_artist_valid_role_check CHECK (role IN ('main', 'featured', 'producer', 'writer')),
    CONSTRAINT unique_track_artist_check UNIQUE (track_id, artist_id, role)
);

-- artist microservice
CREATE TABLE IF NOT EXISTS album_artist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    album_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    FOREIGN KEY (artist_id)
        REFERENCES artist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_album_artist_check UNIQUE (album_id, artist_id)
);

-- playlist microservice
CREATE TABLE IF NOT EXISTS playlist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    CONSTRAINT playlist_title_length_check CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
    description TEXT DEFAULT '',
    CONSTRAINT playlist_description_length_check CHECK (LENGTH(description) <= 1000),
    user_id BIGINT NOT NULL,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    thumbnail_url TEXT NOT NULL DEFAULT '/default_playlist.png',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_user_playlist_check UNIQUE (user_id, title)
);

-- playlist microservice
CREATE TABLE IF NOT EXISTS favorite_playlist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    playlist_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (playlist_id)
        REFERENCES playlist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_user_favorite_playlist_check UNIQUE (user_id, playlist_id)
);

-- playlist microservice
CREATE TABLE IF NOT EXISTS playlist_track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    playlist_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (playlist_id)
        REFERENCES playlist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_playlist_track_check UNIQUE (playlist_id, track_id)
);

-- genre microservice
CREATE TABLE IF NOT EXISTS genre_track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    genre_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (genre_id) 
        REFERENCES genre (id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    CONSTRAINT unique_genre_track_check UNIQUE (genre_id, track_id)
);

-- genre microservice
CREATE TABLE IF NOT EXISTS genre_album (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    genre_id BIGINT NOT NULL,
    album_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (genre_id) 
        REFERENCES genre (id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    CONSTRAINT unique_genre_album_check UNIQUE (genre_id, album_id)
);

-- track microservice
CREATE TABLE IF NOT EXISTS favorite_track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (track_id)
        REFERENCES track (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_favorite_track_check UNIQUE (user_id, track_id)
);

-- album microservice
CREATE TABLE IF NOT EXISTS favorite_album (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    album_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (album_id)
        REFERENCES album (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_favorite_album_check UNIQUE (user_id, album_id)
);

-- artist microservice
CREATE TABLE IF NOT EXISTS favorite_artist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (artist_id)
        REFERENCES artist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_favorite_artist_check UNIQUE (user_id, artist_id)
);

-- track microservice
CREATE TABLE IF NOT EXISTS track_stream (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    duration INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (track_id)
        REFERENCES track (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT stream_valid_duration_check CHECK (duration >= 0)
);

-- track microservice
CREATE MATERIALIZED VIEW track_stats AS
SELECT 
    t.id AS track_id,
    COUNT(DISTINCT ts.user_id) AS listeners_count,
    COUNT(DISTINCT ft.user_id) AS favorites_count
FROM 
    track t
    LEFT JOIN track_stream ts ON t.id = ts.track_id
    LEFT JOIN favorite_track ft ON t.id = ft.track_id
GROUP BY 
    t.id;

CREATE UNIQUE INDEX track_stats_track_id_idx ON track_stats (track_id);

-- album microservice
CREATE TABLE IF NOT EXISTS album_stream (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    album_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (album_id)
        REFERENCES album (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT stream_valid_duration_check CHECK (duration >= 0)
);

-- album microservice
CREATE MATERIALIZED VIEW album_stats AS
SELECT 
    a.id AS album_id,
    COUNT(DISTINCT abs.user_id) AS listeners_count,
    COUNT(DISTINCT fa.user_id) AS favorites_count
FROM 
    album a
    LEFT JOIN album_stream abs ON a.id = abs.album_id
    LEFT JOIN favorite_album fa ON a.id = fa.album_id
GROUP BY 
    a.id;

CREATE UNIQUE INDEX album_stats_album_id_idx ON album_stats (album_id);

-- artist microservice
CREATE TABLE IF NOT EXISTS artist_stream (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (artist_id)
        REFERENCES artist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE MATERIALIZED VIEW artist_stats AS
SELECT 
    a.id AS artist_id,
    COUNT(DISTINCT astr.user_id) AS listeners_count,
    COUNT(DISTINCT fa.user_id) AS favorites_count
FROM 
    artist a
    LEFT JOIN artist_stream astr ON a.id = astr.artist_id
    LEFT JOIN favorite_artist fa ON a.id = fa.artist_id
GROUP BY 
    a.id;

CREATE UNIQUE INDEX artist_stats_artist_id_idx ON artist_stats (artist_id);


SELECT cron.schedule('refresh_artist_stats', '0 * * * *', 'REFRESH MATERIALIZED VIEW CONCURRENTLY artist_stats');
SELECT cron.schedule('refresh_album_stats', '0 * * * *', 'REFRESH MATERIALIZED VIEW CONCURRENTLY album_stats');
SELECT cron.schedule('refresh_track_stats', '0 * * * *', 'REFRESH MATERIALIZED VIEW CONCURRENTLY track_stats');

---- create above / drop below ----

DROP INDEX IF EXISTS artist_stats_artist_id_idx;
DROP INDEX IF EXISTS album_stats_album_id_idx;
DROP INDEX IF EXISTS track_stats_track_id_idx;

DROP MATERIALIZED VIEW IF EXISTS artist_stats;
DROP MATERIALIZED VIEW IF EXISTS album_stats;
DROP MATERIALIZED VIEW IF EXISTS track_stats;

DROP TABLE IF EXISTS track_stream;
DROP TABLE IF EXISTS artist_stream;
DROP TABLE IF EXISTS album_stream;
DROP TABLE IF EXISTS track_stream;
DROP TABLE IF EXISTS favorite_artist;
DROP TABLE IF EXISTS favorite_album;
DROP TABLE IF EXISTS favorite_track;
DROP TABLE IF EXISTS genre_album;
DROP TABLE IF EXISTS genre_track;
DROP TABLE IF EXISTS favorite_playlist;
DROP TABLE IF EXISTS playlist_track;
DROP TABLE IF EXISTS playlist;
DROP TABLE IF EXISTS album_artist;
DROP TABLE IF EXISTS track_artist;
DROP TABLE IF EXISTS track;
DROP TABLE IF EXISTS album;
DROP TABLE IF EXISTS artist;
DROP TABLE IF EXISTS genre;
DROP TABLE IF EXISTS user_settings;
DROP TABLE IF EXISTS "user";

SELECT cron.unschedule('refresh_artist_stats') WHERE EXISTS (SELECT 1 FROM cron.job WHERE jobname = 'refresh_artist_stats');
SELECT cron.unschedule('refresh_album_stats') WHERE EXISTS (SELECT 1 FROM cron.job WHERE jobname = 'refresh_album_stats');
SELECT cron.unschedule('refresh_track_stats') WHERE EXISTS (SELECT 1 FROM cron.job WHERE jobname = 'refresh_track_stats');