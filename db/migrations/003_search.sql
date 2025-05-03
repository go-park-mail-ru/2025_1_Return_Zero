-- Write your migrate up statements here

CREATE EXTENSION IF NOT EXISTS unaccent;

CREATE TEXT SEARCH CONFIGURATION public.multilingual (COPY = pg_catalog.english);
ALTER TEXT SEARCH CONFIGURATION public.multilingual ALTER MAPPING FOR hword, hword_part, word WITH unaccent, russian_stem, english_stem;

ALTER TABLE track ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
  setweight(to_tsvector('english', title), 'A') || 
  setweight(to_tsvector('russian', title), 'A')
) STORED;

ALTER TABLE artist ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
  setweight(to_tsvector('english', title), 'A') || 
  setweight(to_tsvector('russian', title), 'A')
) STORED;

ALTER TABLE album ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
  setweight(to_tsvector('english', title), 'A') || 
  setweight(to_tsvector('russian', title), 'A')
) STORED;

CREATE INDEX track_search_idx ON track USING GIN(search_vector);
CREATE INDEX artist_search_idx ON artist USING GIN(search_vector);
CREATE INDEX album_search_idx ON album USING GIN(search_vector);

---- create above / drop below ----

DROP INDEX IF EXISTS track_search_idx;
DROP INDEX IF EXISTS artist_search_idx;
DROP INDEX IF EXISTS album_search_idx;

ALTER TABLE track DROP COLUMN IF EXISTS search_vector;
ALTER TABLE artist DROP COLUMN IF EXISTS search_vector;
ALTER TABLE album DROP COLUMN IF EXISTS search_vector;

DROP TEXT SEARCH CONFIGURATION IF EXISTS public.multilingual;
DROP EXTENSION IF EXISTS unaccent;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
