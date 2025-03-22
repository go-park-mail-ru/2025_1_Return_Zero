User:
{id} → {email, username, thumbnail_url, password_hash, created_at, updated_at, is_active}

Artist:
{id} → {title, description, created_at, thumbnail_url}

Album:
{id} → {title, type, thumbnail_url, release_date, created_at, artist_id}

Track:
{id} → {title, thumbnail_url, album_id, artist_id, created_at, duration}

Playlist:
{id} → {title, user_id, description, is_public, created_at, updated_at}

Playlist_Track:
{playlist_id, track_id} → {added_at}

Favorite_Track:
{user_id, track_id} → {added_at}

Favorite_Album:
{user_id, album_id} → {added_at}

Favorite_Artist:
{user_id, artist_id} → {added_at}

# 1НФ
- все атрибуты имеют атомарные значения
- нет повторяющихся групп
- у каждой таблицы есть первичный ключ

# 2НФ
- все неключевые атрибуты зависят от первичного ключа
- в составных ключах типа playlist_track дополнительные атрибуты зависят от всей комбинации ключей

# 3НФ
- отсутствуют зависимости неключевых атрибутов от других неключевых атрибутов

# НФБК
- все детерминанты являются потенциальными ключами
