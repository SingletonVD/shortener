ALTER TABLE shortened_links ADD COLUMN IF NOT EXISTS user_id TEXT;

DROP INDEX IF EXISTS idx_shortened_links_full_link;
CREATE UNIQUE INDEX IF NOT EXISTS idx_shortened_links_user_id_full_link ON shortened_links(user_id, full_link);
CREATE INDEX IF NOT EXISTS idx_shortened_links_user_id  ON shortened_links(user_id);
