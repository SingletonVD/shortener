DROP INDEX IF EXISTS idx_shortened_links_user_id;
DROP INDEX IF EXISTS idx_shortened_links_user_id_full_link;
CREATE UNIQUE INDEX IF NOT EXISTS idx_shortened_links_full_link ON shortened_links(full_link);

ALTER TABLE shortened_links DROP COLUMN IF EXISTS user_id;
