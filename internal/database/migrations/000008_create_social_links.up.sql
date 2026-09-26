CREATE TABLE
    social_links (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        type TEXT NOT NULL,
        url TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_social_links_user_id ON social_links (user_id);

CREATE UNIQUE INDEX idx_social_links_user_url ON social_links (user_id, LOWER(url));