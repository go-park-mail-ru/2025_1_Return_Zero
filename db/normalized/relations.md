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

## User
- id - id пользователя
- email - email пользователя
- username - имя пользователя
- thumbnail_url - url аватара пользователя
- password_hash - хэш пароля пользователя
- created_at - дата создания пользователя
- updated_at - дата обновления пользователя
- is_active - активен ли пользователь (пока не будет использоваться, так что на будущее)

## User_Settings
Настройки приватности пользователя 
- user_id - id пользователя
- is_public_playlists - публичные ли плейлисты
- is_public_minutes_listened - публичные ли "минут прослушано"
- is_public_artists_listened - публичные ли "артистов прослушано"
- is_public_favorite_tracks - публичные ли "любимые треки"
- is_public_favorite_albums - публичные ли "любимые альбомы"
- is_public_favorite_artists - публичные ли "любимые исполнители"

## Genre
Жанр музыки
- id - id жанра
- name - название жанра

## Artist
Исполнитель
- id - id исполнителя
- title - название исполнителя
- description - описание исполнителя
- created_at - дата создания исполнителя
- thumbnail_url - url изображения исполнителя

## Album
Альбом
- id - id альбома
- title - название альбома
- type - тип альбома (album, single, ep, compilation)
- thumbnail_url - url изображения альбома
- release_date - дата выпуска альбома
- created_at - дата создания альбома
- artist_id - id исполнителя

## Track
Трек
- id - id трека
- title - название трека
- thumbnail_url - url изображения трека
- album_id - id альбома
- created_at - дата создания трека
- duration - длительность трека
- position - позиция трека в альбоме

## Track_Artist
Связь между треком и исполнителем
- track_id - id трека
- artist_id - id исполнителя
- role - роль исполнителя в треке

## Playlist
Плейлист
- id - id плейлиста
- title - название плейлиста
- user_id - id пользователя
- description - описание плейлиста
- thumbnail_url - url изображения плейлиста
- created_at - дата создания плейлиста
- updated_at - дата обновления плейлиста

## Playlist_Track
Связь между плейлистом и треком
- playlist_id - id плейлиста
- track_id - id трека
- position - позиция трека в плейлисте
- added_at - дата добавления трека в плейлист

## Genre_Track
Связь между жанром и треком
- genre_id - id жанра
- track_id - id трека

## Genre_Album
Связь между жанром и альбомом
- genre_id - id жанра
- album_id - id альбома

## Favorite_Track
Любимый трек
- user_id - id пользователя
- track_id - id трека
- added_at - дата добавления трека в любимые

## Favorite_Album
Любимый альбом
- user_id - id пользователя
- album_id - id альбома
- added_at - дата добавления альбома в любимые

## Favorite_Artist
Любимый исполнитель
- user_id - id пользователя
- artist_id - id исполнителя
- added_at - дата добавления исполнителя в любимые

## Stream
Стрим
- id - id стрима
- user_id - id пользователя
- track_id - id трека
- played_at - дата прослушивания трека
- duration - длительность прослушивания трека (по умолчанию 0, потом через update будем обновлять)
