CREATE TABLE
    experiences (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        role TEXT NOT NULL,
        company TEXT NOT NULL,
        start_date TIMESTAMPTZ NOT NULL,
        end_date TIMESTAMPTZ,
        is_currently_working BOOLEAN NOT NULL DEFAULT FALSE,
        sync_version BIGINT NOT NULL DEFAULT 1,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_experiences_user_id ON experiences (user_id);