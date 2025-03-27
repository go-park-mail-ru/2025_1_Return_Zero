# Функциональные зависимости

User:
{id} → {email, username, thumbnail_url, password_hash, created_at, updated_at, is_active}

User_Settings:
{user_id} → {is_public_playlists, is_public_minutes_listened, is_public_artists_listened, is_public_favorite_tracks, is_public_favorite_albums, is_public_favorite_artists}

Genre:
{id} → {name}

Artist:
{id} → {title, description, created_at, thumbnail_url}

Album:
{id} → {title, type, thumbnail_url, release_date, created_at, artist_id}

Track:
{id} → {title, thumbnail_url, album_id, created_at, duration, position}

Track_Artist:
{track_id, artist_id, role} → {}

Playlist:
{id} → {title, user_id, description, thumbnail_url, created_at, updated_at}

Playlist_Track:
{playlist_id, track_id} → {position, added_at}

Genre_Track:
{genre_id, track_id} → {}

Genre_Album:
{genre_id, album_id} → {}

Genre_Artist:
{genre_id, artist_id} → {}

Favorite_Track:
{user_id, track_id} → {added_at}

Favorite_Album:
{user_id, album_id} → {added_at}

Favorite_Artist:
{user_id, artist_id} → {added_at}

Stream:
{id} → {user_id, track_id, played_at, duration}

## 1НФ
- все атрибуты имеют атомарные значения
- нет повторяющихся групп
- у каждой таблицы есть первичный ключ

## 2НФ
- все неключевые атрибуты зависят от первичного ключа
- в составных ключах типа playlist_track дополнительные атрибуты зависят от всей комбинации ключей

## 3НФ
- отсутствуют зависимости неключевых атрибутов от других неключевых атрибутов

## НФБК
- все детерминанты являются потенциальными ключами

# Описание схемы
