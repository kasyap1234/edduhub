BEGIN;

CREATE TABLE IF NOT EXISTS attendance (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    student_id INT NOT NULL,
    course_id INT NOT NULL,
    college_id INT NOT NULL,
    lecture_id INT NOT NULL,
    date DATE NOT NULL, -- Using DATE type as discussed
    status VARCHAR(50) NOT NULL DEFAULT 'Present',
    scanned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Added created_at/updated_at
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Foreign Keys
    CONSTRAINT fk_attendance_student
        FOREIGN KEY (student_id)
        REFERENCES students(student_id)
        ON DELETE CASCADE,
    CONSTRAINT fk_attendance_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE, -- Or RESTRICT?
    CONSTRAINT fk_attendance_college
        FOREIGN KEY (college_id)
        REFERENCES colleges(id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_attendance_lecture
        FOREIGN KEY (lecture_id)
        REFERENCES lectures(id)
        ON DELETE CASCADE, -- If lecture is deleted, remove attendance?

    -- IMPORTANT: The UNIQUE constraint for your ON CONFLICT clause.
    -- Ensure these columns match your repository code exactly.
    UNIQUE (student_id, course_id, lecture_id, date, college_id)
);

-- Indexes for FKs and common query patterns
CREATE INDEX IF NOT EXISTS idx_attendance_student_id ON attendance (student_id);
CREATE INDEX IF NOT EXISTS idx_attendance_course_id ON attendance (course_id);
CREATE INDEX IF NOT EXISTS idx_attendance_college_id ON attendance (college_id);
CREATE INDEX IF NOT EXISTS idx_attendance_lecture_id ON attendance (lecture_id);
CREATE INDEX IF NOT EXISTS idx_attendance_date ON attendance (date);
-- Index for the unique constraint columns (often created automatically, but explicit doesn't hurt)
CREATE INDEX IF NOT EXISTS idx_attendance_unique_key ON attendance (student_id, course_id, lecture_id, date, college_id);

COMMIT;
