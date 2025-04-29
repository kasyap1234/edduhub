BEGIN;

CREATE TABLE IF NOT EXISTS lectures (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    college_id INT NOT NULL,
    course_id INT NOT NULL,
    -- qr_code_id INT, -- Interpreted as FK to qrcodes(id). Adjust type/name if it refers to the string qr_code_id
    lecture_datetime TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Added a field for when the lecture occurs
    topic VARCHAR(255), -- Added an optional topic field
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_lectures_college
        FOREIGN KEY (college_id)
        REFERENCES colleges(id)
        ON DELETE RESTRICT,
    CONSTRAINT fk_lectures_course
        FOREIGN KEY (course_id)
        REFERENCES courses(id)
        ON DELETE CASCADE -- If course is deleted, maybe delete lectures? Or RESTRICT?
    -- CONSTRAINT fk_lectures_qrcode -- Add if qr_code_id is indeed a FK
    --    FOREIGN KEY (qr_code_id)
    --    REFERENCES qrcodes(id)
    --    ON DELETE SET NULL
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_lectures_college_id ON lectures (college_id);
CREATE INDEX IF NOT EXISTS idx_lectures_course_id ON lectures (course_id);
-- CREATE INDEX IF NOT EXISTS idx_lectures_qr_code_id ON lectures (qr_code_id);
CREATE INDEX IF NOT EXISTS idx_lectures_datetime ON lectures (lecture_datetime);

COMMIT;
