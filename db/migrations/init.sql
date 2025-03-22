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
    CONSTRAINT album_valid_type_check CHECK (type IN ('album', 'single', 'ep', 'compilation'))
);

CREATE TABLE track (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL DEFAULT '/default_track.png',
    album_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    duration INTEGER NOT NULL,
    FOREIGN KEY (album_id)
       REFERENCES album (id)
       ON DELETE CASCADE
       ON UPDATE CASCADE,
    FOREIGN KEY (artist_id)
       REFERENCES artist (id)
       ON DELETE CASCADE
       ON UPDATE CASCADE,
    CONSTRAINT track_valid_duration_check CHECK (duration > 0)
);

CREATE TABLE playlist (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    description TEXT DEFAULT '',
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id)
      REFERENCES user (id)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);

CREATE TABLE playlist_track (
    playlist_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (playlist_id, track_id),
    FOREIGN KEY (playlist_id)
        REFERENCES playlist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (track_id)
        REFERENCES track (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
);

CREATE TABLE favorite_track (
    user_id BIGINT NOT NULL,
    track_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, track_id),
    FOREIGN KEY (user_id)
      REFERENCES user (id)
      ON DELETE CASCADE
      ON UPDATE CASCADE,
    FOREIGN KEY (track_id)
      REFERENCES track (id)
      ON DELETE CASCADE
      ON UPDATE CASCADE
);

CREATE TABLE favorite_album (
    user_id BIGINT NOT NULL,
    album_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, album_id),
    FOREIGN KEY (user_id)
        REFERENCES user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (album_id)
        REFERENCES album (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE favorite_artist (
    user_id BIGINT NOT NULL,
    artist_id BIGINT NOT NULL,
    added_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, artist_id),
    FOREIGN KEY (user_id)
        REFERENCES user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (artist_id)
        REFERENCES artist (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);