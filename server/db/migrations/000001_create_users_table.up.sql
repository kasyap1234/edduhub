BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    -- Assuming you link users via Kratos ID
    kratos_identity_id VARCHAR(255) NOT NULL UNIQUE,
    -- Add other fields you might need, like email or role
    -- email VARCHAR(255) UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;