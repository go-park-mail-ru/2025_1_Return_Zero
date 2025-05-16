CREATE INDEX idx_playlist_user_created_at ON playlist(user_id, created_at DESC);
CREATE INDEX idx_playlist_public ON playlist(is_public);
CREATE INDEX idx_favorite_playlist_user_playlist ON favorite_playlist(user_id, playlist_id, created_at DESC);
