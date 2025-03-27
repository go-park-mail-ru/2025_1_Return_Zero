```mermaid
erDiagram
    user {
        BIGINT id PK
        TEXT email
        TEXT username
        TEXT thumbnail_url
        TEXT password_hash
        TIMESTAMP created_at
        TIMESTAMP updated_at
        BOOLEAN is_active
    }

    user_settings {
        BIGINT user_id PK
        BOOLEAN is_public_playlists
        BOOLEAN is_public_minutes_listened
        BOOLEAN is_public_artists_listened
        BOOLEAN is_public_favorite_tracks
        BOOLEAN is_public_favorite_albums
        BOOLEAN is_public_favorite_artists
    }

    genre {
        BIGINT id PK
        TEXT name
    }

    artist {
        BIGINT id PK
        TEXT title
        TEXT description
        TIMESTAMP created_at
        TEXT thumbnail_url
    }

    album {
        BIGINT id PK
        TEXT title
        TEXT type
        TEXT thumbnail_url
        DATE release_date
        TIMESTAMP created_at
        BIGINT artist_id FK
    }

    track {
        BIGINT id PK
        TEXT title
        TEXT thumbnail_url
        BIGINT album_id FK
        TIMESTAMP created_at
        INTEGER duration
        INTEGER position
    }

    track_artist {
        BIGINT id PK
        BIGINT track_id FK
        BIGINT artist_id FK
        TEXT role
    }

    playlist {
        BIGINT id PK
        TEXT title
        TEXT description
        BIGINT user_id FK
        TEXT thumbnail_url
        TIMESTAMP created_at
        TIMESTAMP updated_at
    }

    playlist_track {
        BIGINT id PK
        BIGINT playlist_id FK
        BIGINT track_id FK
        BIGINT position
        TIMESTAMP added_at
    }

    genre_track {
        BIGINT id PK
        BIGINT genre_id FK
        BIGINT track_id FK
    }

    genre_album {
        BIGINT id PK
        BIGINT genre_id FK
        BIGINT album_id FK
    }

    favorite_track {
        BIGINT id PK
        BIGINT user_id FK
        BIGINT track_id FK
        TIMESTAMP added_at
    }

    favorite_album {
        BIGINT id PK
        BIGINT user_id FK
        BIGINT album_id FK
        TIMESTAMP added_at
    }

    favorite_artist {
        BIGINT id PK
        BIGINT user_id FK
        BIGINT artist_id FK
        TIMESTAMP added_at
    }

    stream {
        BIGINT id PK
        BIGINT user_id FK
        BIGINT track_id FK
        TIMESTAMP played_at
        INTEGER duration
    }

    user ||--|| user_settings : "has"
    user ||--o{ playlist : "creates"
    user ||--o{ stream : "listens to"
    user ||--o{ favorite_track : "favorites"
    user ||--o{ favorite_album : "favorites"
    user ||--o{ favorite_artist : "favorites"
    
    artist ||--o{ album : "has"
    artist ||--o{ track_artist : "contributes to"
    artist ||--o{ favorite_artist : "favorited_by"
    
    album ||--o{ track : "contains"
    album ||--o{ genre_album : "categorized as"
    album ||--o{ favorite_album : "favorited_by"
    
    track ||--o{ track_artist : "has"
    track ||--o{ playlist_track : "included_in"
    track ||--o{ genre_track : "categorized as"
    track ||--o{ favorite_track : "favorited_by"
    track ||--o{ stream : "streamed as"
    
    playlist ||--o{ playlist_track : "includes"
    
    genre ||--o{ genre_track : "categorizes"
    genre ||--o{ genre_album : "categorizes"
```
