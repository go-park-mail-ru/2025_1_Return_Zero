# Функциональные зависимости

User:
{id} → {email, username, thumbnail_url, password_hash, created_at, updated_at, is_active}

User_Settings:
{user_id} → {is_public_playlists, is_public_minutes_listened, is_public_artists_listened, is_public_favorite_tracks, is_public_favorite_albums, is_public_favorite_artists}

Genre:
{id} → {name}

Artist:
{id} → {title, description, created_at, thumbnail_url, listeners_count, favorites_count}

Album:
{id} → {title, type, thumbnail_url, release_date, created_at, artist_id, listeners_count, favorites_count}

Track:
{id} → {title, thumbnail_url, file_url, album_id, created_at, duration, position, listeners_count, favorites_count}

Track_Artist:
{id} → {track_id, artist_id, role}

Playlist:
{id} → {title, description, user_id, thumbnail_url, created_at, updated_at}

Playlist_Track:
{id} → {playlist_id, track_id, position, added_at}

Genre_Track:
{id} → {genre_id, track_id}

Genre_Album:
{id} → {genre_id, album_id}

Favorite_Track:
{id} → {user_id, track_id, added_at}

Favorite_Album:
{id} → {user_id, album_id, added_at}

Favorite_Artist:
{id} → {user_id, artist_id, added_at}

Stream:
{id} → {user_id, track_id, played_at, duration}

## 1НФ
- все атрибуты имеют атомарные значения
- нет повторяющихся групп
- у каждой таблицы есть первичный ключ

## 2НФ
- все неключевые атрибуты зависят от первичного ключа

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
- listeners_count - количество слушателей
- favorites_count - количество добавлений в избранное
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
- listeners_count - количество слушателей
- favorites_count - количество добавлений в избранное

## Track
Трек
- id - id трека
- title - название трека
- thumbnail_url - url изображения трека
- file_url - url файла трека
- album_id - id альбома
- created_at - дата создания трека
- duration - длительность трека
- position - позиция трека в альбоме
- listeners_count - количество слушателей
- favorites_count - количество добавлений в избранное

## Track_Artist
Связь между треком и исполнителем
- id - id связи
- track_id - id трека
- artist_id - id исполнителя
- role - роль исполнителя в треке (main, featured, producer, writer)

## Playlist
Плейлист
- id - id плейлиста
- title - название плейлиста
- description - описание плейлиста
- user_id - id пользователя (может быть NULL для системных подборок)
- thumbnail_url - url изображения плейлиста
- created_at - дата создания плейлиста
- updated_at - дата обновления плейлиста

## Playlist_Track
Связь между плейлистом и треком
- id - id связи
- playlist_id - id плейлиста
- track_id - id трека
- position - позиция трека в плейлисте
- added_at - дата добавления трека в плейлист

## Genre_Track
Связь между жанром и треком
- id - id связи
- genre_id - id жанра
- track_id - id трека

## Genre_Album
Связь между жанром и альбомом
- id - id связи
- genre_id - id жанра
- album_id - id альбома

## Favorite_Track
Любимый трек
- id - id записи
- user_id - id пользователя
- track_id - id трека
- added_at - дата добавления трека в любимые

## Favorite_Album
Любимый альбом
- id - id записи
- user_id - id пользователя
- album_id - id альбома
- added_at - дата добавления альбома в любимые

## Favorite_Artist
Любимый исполнитель
- id - id записи
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
