-- Write your migrate up statements here

DROP INDEX IF EXISTS track_stats_track_id_idx;
DROP MATERIALIZED VIEW IF EXISTS track_stats;

CREATE MATERIALIZED VIEW track_stats AS
SELECT 
    t.id AS track_id,
    COUNT(DISTINCT ts.user_id) AS listeners_count,
    COUNT(DISTINCT ft.user_id) AS favorites_count,
    COUNT(DISTINCT CASE 
        WHEN ts.created_at >= NOW() - INTERVAL '1 month' 
        THEN ts.user_id 
        ELSE NULL 
    END) AS listeners_count_last_month,
    COUNT(DISTINCT CASE 
        WHEN ft.created_at >= NOW() - INTERVAL '1 week' 
        THEN ft.user_id 
        ELSE NULL 
    END) AS favorites_count_last_week
FROM 
    track t
    LEFT JOIN track_stream ts ON t.id = ts.track_id
    LEFT JOIN favorite_track ft ON t.id = ft.track_id
GROUP BY 
    t.id;

CREATE UNIQUE INDEX track_stats_track_id_idx ON track_stats (track_id);

REFRESH MATERIALIZED VIEW track_stats;

---- create above / drop below ----
DROP INDEX IF EXISTS track_stats_track_id_idx;
DROP MATERIALIZED VIEW IF EXISTS track_stats;

CREATE MATERIALIZED VIEW track_stats AS
SELECT 
    t.id AS track_id,
    COUNT(DISTINCT ts.user_id) AS listeners_count,
    COUNT(DISTINCT ft.user_id) AS favorites_count
FROM 
    track t
    LEFT JOIN track_stream ts ON t.id = ts.track_id
    LEFT JOIN favorite_track ft ON t.id = ft.track_id
GROUP BY 
    t.id;

CREATE UNIQUE INDEX track_stats_track_id_idx ON track_stats (track_id);

REFRESH MATERIALIZED VIEW track_stats;
