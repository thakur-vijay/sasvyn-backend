CREATE TABLE
    experience_responsibilities (
        id UUID PRIMARY KEY,
        experience_id UUID NOT NULL REFERENCES experiences (id) ON DELETE CASCADE,
        responsibility TEXT NOT NULL,
        "order" INTEGER NOT NULL,
        sync_version BIGINT NOT NULL DEFAULT 1,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_experience_responsibilities_experience_id ON experience_responsibilities (experience_id);