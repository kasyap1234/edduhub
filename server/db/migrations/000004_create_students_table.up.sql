BEGIN;

CREATE TABLE IF NOT EXISTS students (
    -- Assuming student_id is the primary key, matching the model tag
    student_id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id INT NOT NULL UNIQUE, -- Each user can only be one student
    college_id INT NOT NULL,
    kratos_identity_id VARCHAR(255) NOT NULL UNIQUE, -- Matches the user's kratos ID
    enrollment_year INT,
    roll_no VARCHAR(50) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign Keys
    CONSTRAINT fk_students_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE, -- If user is deleted, delete the student record
    CONSTRAINT fk_students_college
        FOREIGN KEY (college_id)
        REFERENCES colleges(id)
        ON DELETE RESTRICT, -- Don't allow deleting a college if students are linked
    CONSTRAINT fk_students_kratos -- Ensure student kratos ID matches user kratos ID
        FOREIGN KEY (kratos_identity_id)
        REFERENCES users(kratos_identity_id)
        ON DELETE CASCADE, -- If user is deleted

    -- Unique constraint for roll number within a college
    UNIQUE (college_id, roll_no)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_students_user_id ON students (user_id);
CREATE INDEX IF NOT EXISTS idx_students_college_id ON students (college_id);
CREATE INDEX IF NOT EXISTS idx_students_kratos_identity_id ON students (kratos_identity_id);
CREATE INDEX IF NOT EXISTS idx_students_roll_no ON students (college_id, roll_no); -- For the unique constraint

COMMIT;
