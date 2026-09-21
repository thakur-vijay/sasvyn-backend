CREATE TABLE
    skills (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        skill TEXT NOT NULL,
        category TEXT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_skills_user_id ON skills (user_id);

CREATE UNIQUE INDEX idx_skills_user_skill ON skills (user_id, LOWER(skill));