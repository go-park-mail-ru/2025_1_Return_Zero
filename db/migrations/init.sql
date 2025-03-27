CREATE TABLE user (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    thumbnail_url TEXT NOT NULL DEFAULT '/default_avatar.png',
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    CONSTRAINT user_valid_email_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
);

CREATE TABLE user_settings (
    user_id BIGINT NOT NULL PRIMARY KEY,
    is_public_playlists BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_minutes_listened BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_artists_listened BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_favorite_tracks BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_favorite_albums BOOLEAN NOT NULL DEFAULT FALSE,
    is_public_favorite_artists BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (user_id)
        REFERENCES user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE genre (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE artist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    thumbnail_url TEXT NOT NULL DEFAULT '/default_artist.png'
);

CREATE TABLE album (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    type TEXT NOT NULL DEFAULT 'single',
    thumbnail_url TEXT NOT NULL DEFAULT '/default_album.png',
    release_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    artist_id BIGINT NOT NULL,
    FOREIGN KEY (artist_id)
        REFERENCES artist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT album_valid_type_check CHECK (type IN ('album', 'single', 'ep', 'compilation')),
    CONSTRAINT unique_artist_album_check UNIQUE (artist_id, title)
);

CREATE TABLE track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL DEFAULT '/default_track.png',
    album_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    duration INTEGER NOT NULL,
    position INTEGER NOT NULL,
    FOREIGN KEY (album_id)
        REFERENCES album (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT track_valid_duration_check CHECK (duration > 0),
    CONSTRAINT unique_album_track_check UNIQUE (album_id, position)
);

CREATE TABLE track_artist (
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

CREATE TABLE playlist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    user_id BIGINT, -- NULL для подборок
    thumbnail_url TEXT NOT NULL DEFAULT '/default_playlist.png',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
        REFERENCES user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_user_playlist_check UNIQUE (user_id, title)
);

CREATE TABLE playlist_track (
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

CREATE TABLE genre_track (
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

CREATE TABLE genre_album (
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

CREATE TABLE favorite_track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
        REFERENCES user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (track_id)
        REFERENCES track (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_favorite_track_check UNIQUE (user_id, track_id)
);

CREATE TABLE favorite_album (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    album_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
        REFERENCES user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (album_id)
        REFERENCES album (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_favorite_album_check UNIQUE (user_id, album_id)
);

CREATE TABLE favorite_artist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
        REFERENCES user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (artist_id)
        REFERENCES artist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT unique_favorite_artist_check UNIQUE (user_id, artist_id)
);

CREATE TABLE stream (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    played_at TIMESTAMP NOT NULL DEFAULT NOW(),
    duration INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (user_id)
        REFERENCES user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (track_id)
        REFERENCES track (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT stream_valid_duration_check CHECK (duration >= 0)
);
