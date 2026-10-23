-- Create lesson_codes table (master data per SPEC.md §3.2/§5)
CREATE TABLE lesson_codes (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    type VARCHAR(100) NOT NULL,
    default_fee BIGINT,
    duration_hours INT,
    description TEXT,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_lesson_codes_active ON lesson_codes(active);
