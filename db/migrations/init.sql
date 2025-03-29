-- Write your migrate up statements here

CREATE TABLE IF NOT EXISTS "user" (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    CONSTRAINT email_length_check CHECK (LENGTH(email) >= 5 AND LENGTH(email) <= 30),
    CONSTRAINT user_valid_email_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    username TEXT NOT NULL UNIQUE,
    CONSTRAINT username_length_check CHECK (LENGTH(username) >= 3 AND LENGTH(username) <= 20),
    thumbnail_url TEXT NOT NULL DEFAULT '/default_avatar.png',
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE IF NOT EXISTS user_settings (
    user_id BIGINT NOT NULL PRIMARY KEY,
    is_public_playlists BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_minutes_listened BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_artists_listened BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_favorite_tracks BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_favorite_albums BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_favorite_artists BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (user_id)
        REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS genre (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    CONSTRAINT genre_name_length_check CHECK (LENGTH(name) >= 3 AND LENGTH(name) <= 20)
);

CREATE TABLE IF NOT EXISTS artist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    CONSTRAINT artist_title_length_check CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
    description TEXT NOT NULL DEFAULT '',
    CONSTRAINT artist_description_length_check CHECK (LENGTH(description) <= 1000),
    listeners_count BIGINT NOT NULL DEFAULT 0,
    favorites_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    thumbnail_url TEXT NOT NULL DEFAULT '/default_artist.png',
    CONSTRAINT non_negative_listeners_count_check CHECK (listeners_count >= 0),
    CONSTRAINT non_negative_favorites_count_check CHECK (favorites_count >= 0)
);

CREATE TABLE IF NOT EXISTS album (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    CONSTRAINT album_title_length_check CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
    type TEXT NOT NULL DEFAULT 'single',
    thumbnail_url TEXT NOT NULL DEFAULT '/default_album.png',
    release_date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    artist_id BIGINT NOT NULL,
    listeners_count BIGINT NOT NULL DEFAULT 0,
    favorites_count BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (artist_id)
        REFERENCES artist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT album_valid_type_check CHECK (type IN ('album', 'single', 'ep', 'compilation')),
    CONSTRAINT unique_artist_album_check UNIQUE (artist_id, title),
    CONSTRAINT non_negative_listeners_count_check CHECK (listeners_count >= 0),
    CONSTRAINT non_negative_favorites_count_check CHECK (favorites_count >= 0)
);

CREATE TABLE IF NOT EXISTS track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    CONSTRAINT track_title_length_check CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
    thumbnail_url TEXT NOT NULL DEFAULT '/default_track.png',
    file_url TEXT NOT NULL DEFAULT '',
    album_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    duration INTEGER NOT NULL,
    position INTEGER NOT NULL,
    listeners_count BIGINT NOT NULL DEFAULT 0,
    favorites_count BIGINT NOT NULL DEFAULT 0,
    FOREIGN KEY (album_id)
        REFERENCES album (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT track_valid_duration_check CHECK (duration > 0),
    CONSTRAINT unique_album_track_check UNIQUE (album_id, position),
    CONSTRAINT non_negative_listeners_count_check CHECK (listeners_count >= 0),
    CONSTRAINT non_negative_favorites_count_check CHECK (favorites_count >= 0)
);

CREATE TABLE IF NOT EXISTS track_artist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, 
    track_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    role TEXT NOT NULL DEFAULT 'main',
    FOREIGN KEY (track_id) 
        REFERENCES track (id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    FOREIGN KEY (artist_id) 
        REFERENCES artist (id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    CONSTRAINT track_artist_valid_role_check CHECK (role IN ('main', 'featured', 'producer', 'writer')),
    CONSTRAINT unique_track_artist_check UNIQUE (track_id, artist_id, role)
);

CREATE TABLE IF NOT EXISTS playlist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    CONSTRAINT playlist_title_length_check CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 100),
    description TEXT DEFAULT '',
    CONSTRAINT playlist_description_length_check CHECK (LENGTH(description) <= 1000),
    user_id BIGINT, -- NULL для подборок
    thumbnail_url TEXT NOT NULL DEFAULT '/default_playlist.png',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
        REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_user_playlist_check UNIQUE (user_id, title)
);

CREATE TABLE IF NOT EXISTS playlist_track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    playlist_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    position BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (playlist_id)
        REFERENCES playlist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (track_id)
        REFERENCES track (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_playlist_position_check UNIQUE (playlist_id, position),
    CONSTRAINT unique_playlist_track_check UNIQUE (playlist_id, track_id)
);

CREATE TABLE IF NOT EXISTS genre_track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    genre_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    FOREIGN KEY (genre_id) 
        REFERENCES genre (id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    FOREIGN KEY (track_id) 
        REFERENCES track (id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    CONSTRAINT unique_genre_track_check UNIQUE (genre_id, track_id)
);

CREATE TABLE IF NOT EXISTS genre_album (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    genre_id BIGINT NOT NULL,
    album_id BIGINT NOT NULL,
    FOREIGN KEY (genre_id) 
        REFERENCES genre (id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    FOREIGN KEY (album_id) 
        REFERENCES album (id) 
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    CONSTRAINT unique_genre_album_check UNIQUE (genre_id, album_id)
);

CREATE TABLE IF NOT EXISTS favorite_track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
        REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (track_id)
        REFERENCES track (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_favorite_track_check UNIQUE (user_id, track_id)
);

CREATE TABLE IF NOT EXISTS favorite_album (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    album_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
        REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (album_id)
        REFERENCES album (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_favorite_album_check UNIQUE (user_id, album_id)
);

CREATE TABLE IF NOT EXISTS favorite_artist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
        REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (artist_id)
        REFERENCES artist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_favorite_artist_check UNIQUE (user_id, artist_id)
);

CREATE TABLE IF NOT EXISTS stream (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    played_at TIMESTAMP NOT NULL DEFAULT NOW(),
    duration INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id)
        REFERENCES "user" (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (track_id)
        REFERENCES track (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT stream_valid_duration_check CHECK (duration >= 0)
);

---- create above / drop below ----

DROP TABLE IF EXISTS stream;
DROP TABLE IF EXISTS favorite_artist;
DROP TABLE IF EXISTS favorite_album;
DROP TABLE IF EXISTS favorite_track;
DROP TABLE IF EXISTS genre_album;
DROP TABLE IF EXISTS genre_track;
DROP TABLE IF EXISTS playlist_track;
DROP TABLE IF EXISTS playlist;
DROP TABLE IF EXISTS track_artist;
DROP TABLE IF EXISTS track;
DROP TABLE IF EXISTS album;
DROP TABLE IF EXISTS artist;
DROP TABLE IF EXISTS genre;
DROP TABLE IF EXISTS user_settings;
DROP TABLE IF EXISTS "user";
