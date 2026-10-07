CREATE TABLE
    documents (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        name TEXT NOT NULL,
        category TEXT NOT NULL,
        file_size BIGINT NOT NULL,
        key TEXT NOT NULL,
        sync_version BIGINT NOT NULL DEFAULT 1,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_documents_user_id ON documents (user_id);