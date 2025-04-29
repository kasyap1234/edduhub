BEGIN;

CREATE TABLE IF NOT EXISTS colleges (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    city VARCHAR(100),
    state VARCHAR(100),
    country VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for potential lookups by name
CREATE INDEX IF NOT EXISTS idx_colleges_name ON colleges (name);

COMMIT;
