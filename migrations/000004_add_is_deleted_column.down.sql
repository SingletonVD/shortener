DROP INDEX IF EXISTS idx_shortened_links_full_link;
CREATE UNIQUE INDEX IF NOT EXISTS idx_shortened_links_user_id_full_link ON shortened_links(user_id, full_link);

DROP INDEX IF EXISTS idx_shortened_links_user_id;
CREATE INDEX IF NOT EXISTS idx_shortened_links_user_id ON shortened_links(user_id);


ALTER TABLE shortened_links DROP COLUMN IF EXISTS is_deleted;
