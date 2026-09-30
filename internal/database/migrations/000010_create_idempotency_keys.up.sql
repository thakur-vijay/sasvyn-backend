CREATE TABLE
    idempotency_keys (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
        idempotency_key TEXT NOT NULL,
        status_code INT NOT NULL,
        response_body JSONB NOT NULL,
        created_at TIMESTAMPTZ NOT NULL DEFAULT NOW (),
        UNIQUE (user_id, idempotency_key)
    );