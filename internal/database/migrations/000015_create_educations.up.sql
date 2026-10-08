CREATE TABLE
    educations (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        degree TEXT NOT NULL,
        field_of_study TEXT NOT NULL,
        institution TEXT NOT NULL,
        start_date TIMESTAMPTZ NOT NULL,
        end_date TIMESTAMPTZ,
        is_pursuing BOOLEAN NOT NULL DEFAULT FALSE,
        grade TEXT NOT NULL,
        grade_type TEXT NOT NULL,
        description TEXT,
        sync_version BIGINT NOT NULL DEFAULT 1,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_educations_user_id ON educations (user_id);