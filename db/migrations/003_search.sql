-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE TEXT SEARCH CONFIGURATION public.multilingual (COPY = pg_catalog.english);
ALTER TEXT SEARCH CONFIGURATION public.multilingual ALTER MAPPING FOR hword, hword_part, word WITH unaccent, russian_stem, english_stem;

ALTER TABLE track ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
  setweight(to_tsvector('english', title), 'A') || 
  setweight(to_tsvector('russian', title), 'A')
) STORED;

ALTER TABLE track ADD COLUMN title_trgm text 
GENERATED ALWAYS AS (title) STORED;

ALTER TABLE artist ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
  setweight(to_tsvector('english', title), 'A') || 
  setweight(to_tsvector('russian', title), 'A')
) STORED;

ALTER TABLE artist ADD COLUMN title_trgm text
GENERATED ALWAYS AS (title) STORED;

ALTER TABLE album ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
  setweight(to_tsvector('english', title), 'A') || 
  setweight(to_tsvector('russian', title), 'A')
) STORED;

ALTER TABLE album ADD COLUMN title_trgm text
GENERATED ALWAYS AS (title) STORED;

ALTER TABLE playlist ADD COLUMN search_vector tsvector
GENERATED ALWAYS AS (
  setweight(to_tsvector('english', title), 'A') || 
  setweight(to_tsvector('russian', title), 'A')
) STORED;

ALTER TABLE playlist ADD COLUMN title_trgm text
GENERATED ALWAYS AS (title) STORED;

CREATE INDEX track_search_idx ON track USING GIN(search_vector);
CREATE INDEX artist_search_idx ON artist USING GIN(search_vector);
CREATE INDEX album_search_idx ON album USING GIN(search_vector);
CREATE INDEX playlist_search_idx ON playlist USING GIN(search_vector);

CREATE INDEX track_trgm_idx ON track USING GIN(title_trgm gin_trgm_ops);
CREATE INDEX artist_trgm_idx ON artist USING GIN(title_trgm gin_trgm_ops);
CREATE INDEX album_trgm_idx ON album USING GIN(title_trgm gin_trgm_ops);
CREATE INDEX playlist_trgm_idx ON playlist USING GIN(title_trgm gin_trgm_ops);
---- create above / drop below ----

DROP INDEX IF EXISTS track_trgm_idx;
DROP INDEX IF EXISTS artist_trgm_idx;
DROP INDEX IF EXISTS album_trgm_idx;

DROP INDEX IF EXISTS track_search_idx;
DROP INDEX IF EXISTS artist_search_idx;
DROP INDEX IF EXISTS album_search_idx;

ALTER TABLE track DROP COLUMN IF EXISTS title_trgm;
ALTER TABLE artist DROP COLUMN IF EXISTS title_trgm;
ALTER TABLE album DROP COLUMN IF EXISTS title_trgm;
ALTER TABLE playlist DROP COLUMN IF EXISTS title_trgm;

ALTER TABLE track DROP COLUMN IF EXISTS search_vector;
ALTER TABLE artist DROP COLUMN IF EXISTS search_vector;
ALTER TABLE album DROP COLUMN IF EXISTS search_vector;
ALTER TABLE playlist DROP COLUMN IF EXISTS search_vector;

DROP TEXT SEARCH CONFIGURATION IF EXISTS public.multilingual;
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS unaccent;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
