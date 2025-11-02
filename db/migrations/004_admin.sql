INSERT INTO "user" (username, password_hash, email, thumbnail_url, created_at, updated_at)
VALUES ('admin', 'AQIDBAUGBwjZy6L5y+0MZhP2MUf9TYynMMdaJmjVm766KGjOeODhDQ==', 'admin@admin.ru', '/default_avatar.png', NOW(), NOW());

INSERT INTO "user_settings" (user_id, is_public_playlists, is_public_minutes_listened, is_public_favorite_artists, is_public_tracks_listened, is_public_favorite_tracks, is_public_artists_listened)
VALUES (
  (SELECT id FROM "user" WHERE username = 'admin'),
  false, false, false, false, false, false
);

---- create above / drop below ----

DELETE FROM "user_settings"
WHERE user_id = (SELECT id FROM "user" WHERE username = 'admin');

DELETE FROM "user"
WHERE username = 'admin';
