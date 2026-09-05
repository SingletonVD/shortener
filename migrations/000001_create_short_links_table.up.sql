CREATE TABLE IF NOT EXISTS shortened_links (
    id SERIAL PRIMARY KEY,
    short_link VARCHAR(8) NOT NULL,
    full_link TEXT NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_shortened_links_short_link ON shortened_links(short_link);
