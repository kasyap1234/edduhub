BEGIN;

CREATE TABLE IF NOT EXISTS courses (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    credits INT NOT NULL DEFAULT 0,
    instructor_id INT, -- Nullable if instructor can be unassigned initially
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign key to users table (assuming instructors are users)
    CONSTRAINT fk_courses_instructor
        FOREIGN KEY (instructor_id)
        REFERENCES users(id)
        ON DELETE SET NULL -- Or RESTRICT / NO ACTION if instructor deletion shouldn't be allowed
);

-- Index for foreign key and potential lookups by name
CREATE INDEX IF NOT EXISTS idx_courses_instructor_id ON courses (instructor_id);
CREATE INDEX IF NOT EXISTS idx_courses_name ON courses (name);

COMMIT;
