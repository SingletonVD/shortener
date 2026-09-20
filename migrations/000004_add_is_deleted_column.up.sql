ALTER TABLE shortened_links ADD COLUMN IF NOT EXISTS is_deleted BOOLEAN NOT NULL DEFAULT FALSE;

DROP INDEX IF EXISTS idx_shortened_links_user_id_full_link;
CREATE UNIQUE INDEX IF NOT EXISTS idx_shortened_links_user_id_full_link ON shortened_links(user_id, full_link)
    WHERE is_deleted = FALSE;

DROP INDEX IF EXISTS idx_shortened_links_user_id;
CREATE INDEX IF NOT EXISTS idx_shortened_links_user_id ON shortened_links(user_id) WHERE is_deleted = FALSE;
