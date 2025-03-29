CREATE TABLE IF NOT EXISTS user (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    CONSTRAINT email_length_check CHECK (LENGTH(email) >= 5 AND LENGTH(email) <= 30),
    CONSTRAINT user_valid_email_check CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$')
    username TEXT NOT NULL UNIQUE,
    CONSTRAINT username_length_check CHECK (LENGTH(username) >= 3 AND LENGTH(username) <= 20),
    thumbnail_url TEXT NOT NULL DEFAULT '/default_avatar.png',
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
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
        REFERENCES user (id)
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
        REFERENCES user (id)
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
        REFERENCES user (id)
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
        REFERENCES user (id)
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
        REFERENCES user (id)
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
        REFERENCES user (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (track_id)
        REFERENCES track (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT stream_valid_duration_check CHECK (duration >= 0)
);

CREATE OR REPLACE FUNCTION update_track_listeners_count()
RETURNS TRIGGER AS $$
DECLARE
    target_id BIGINT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        target_id := NEW.track_id;
        UPDATE track SET listeners_count = listeners_count + 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    ELSIF TG_OP = 'DELETE' THEN
        target_id := OLD.track_id;
        UPDATE track SET listeners_count = listeners_count - 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_track_favorites_count()
RETURNS TRIGGER AS $$
DECLARE
    target_id BIGINT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        target_id := NEW.track_id;
        UPDATE track SET favorites_count = favorites_count + 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    ELSIF TG_OP = 'DELETE' THEN
        target_id := OLD.track_id;
        UPDATE track SET favorites_count = favorites_count - 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_album_favorites_count()
RETURNS TRIGGER AS $$
DECLARE
    target_id BIGINT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        target_id := NEW.album_id;
        UPDATE album SET favorites_count = favorites_count + 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    ELSIF TG_OP = 'DELETE' THEN
        target_id := OLD.album_id;
        UPDATE album SET favorites_count = favorites_count - 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_album_listeners_count()
RETURNS TRIGGER AS $$
DECLARE
    target_id BIGINT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        SELECT album_id INTO target_id FROM track WHERE id = NEW.track_id FOR UPDATE;
        UPDATE album SET listeners_count = listeners_count + 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    ELSIF TG_OP = 'DELETE' THEN
        SELECT album_id INTO target_id FROM track WHERE id = OLD.track_id FOR UPDATE;
        UPDATE album SET listeners_count = listeners_count - 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_artist_favorites_count()
RETURNS TRIGGER AS $$
DECLARE
    target_id BIGINT;
BEGIN
    IF TG_OP = 'INSERT' THEN
        target_id := NEW.artist_id;
        UPDATE artist SET favorites_count = favorites_count + 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    ELSIF TG_OP = 'DELETE' THEN
        target_id := OLD.artist_id;
        UPDATE artist SET favorites_count = favorites_count - 1 WHERE id = target_id AND pg_try_advisory_xact_lock(target_id);
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION update_artist_listeners_count()
RETURNS TRIGGER AS $$
DECLARE
    artist_ids BIGINT[];
BEGIN
    IF TG_OP = 'INSERT' THEN
        SELECT array_agg(DISTINCT ta.artist_id) INTO artist_ids
        FROM track_artist ta
        WHERE ta.track_id = NEW.track_id AND (ta.role = 'main' OR ta.role = 'featured')
        FOR UPDATE;
        
        IF artist_ids IS NOT NULL THEN
            FOREACH target_id IN ARRAY artist_ids LOOP
                PERFORM pg_advisory_xact_lock(target_id);
                UPDATE artist
                SET listeners_count = listeners_count + 1
                WHERE id = target_id;
            END LOOP;
        END IF;
    ELSIF TG_OP = 'DELETE' THEN
        SELECT array_agg(DISTINCT ta.artist_id) INTO artist_ids
        FROM track_artist ta
        WHERE ta.track_id = OLD.track_id AND (ta.role = 'main' OR ta.role = 'featured')
        FOR UPDATE;
        
        IF artist_ids IS NOT NULL THEN
            FOREACH target_id IN ARRAY artist_ids LOOP
                PERFORM pg_advisory_xact_lock(target_id);
                UPDATE artist
                SET listeners_count = listeners_count - 1
                WHERE id = target_id;
            END LOOP;
        END IF;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_track_listeners_count_trigger
    AFTER INSERT OR DELETE ON stream
    FOR EACH ROW
    EXECUTE FUNCTION update_track_listeners_count();

CREATE TRIGGER update_album_listeners_count_trigger
    AFTER INSERT OR DELETE ON stream
    FOR EACH ROW
    EXECUTE FUNCTION update_album_listeners_count();

CREATE TRIGGER update_track_favorites_count_trigger
    AFTER INSERT OR DELETE ON favorite_track
    FOR EACH ROW
    EXECUTE FUNCTION update_track_favorites_count();

CREATE TRIGGER update_album_favorites_count_trigger
    AFTER INSERT OR DELETE ON favorite_album
    FOR EACH ROW
    EXECUTE FUNCTION update_album_favorites_count();

CREATE TRIGGER update_artist_favorites_count_trigger
    AFTER INSERT OR DELETE ON favorite_artist
    FOR EACH ROW
    EXECUTE FUNCTION update_artist_favorites_count();

CREATE TRIGGER update_artist_listeners_count_trigger
    AFTER INSERT OR DELETE ON stream
    FOR EACH ROW
    EXECUTE FUNCTION update_artist_listeners_count(); 