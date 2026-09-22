CREATE TABLE
    languages (
        id UUID PRIMARY KEY,
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        language_code TEXT NOT NULL,
        language TEXT NOT NULL,
        proficiency SMALLINT NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW ()
    );

CREATE INDEX idx_languages_user_id ON languages (user_id);

CREATE UNIQUE INDEX idx_languages_user_language ON languages (user_id, LOWER(language));