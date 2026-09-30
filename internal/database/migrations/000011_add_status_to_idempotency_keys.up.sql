ALTER TABLE idempotency_keys
ALTER COLUMN status_code
DROP NOT NULL,
ALTER COLUMN response_body
DROP NOT NULL;

ALTER TABLE idempotency_keys
ADD COLUMN status TEXT NOT NULL DEFAULT 'completed';