BEGIN;

CREATE TABLE IF NOT EXISTS users (
       id                 INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    kratos_identity_id VARCHAR(255) NOT NULL UNIQUE, -- Link to Kratos
    name               VARCHAR(255) NOT NULL,        -- User's full name
    role               VARCHAR(50)  NOT NULL,        -- e.g., 'student', 'admin', 'instructor'
    email              VARCHAR(255) NOT NULL UNIQUE, -- User's email, should be unique
    is_active          BOOLEAN      NOT NULL DEFAULT TRUE, -- User status
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;