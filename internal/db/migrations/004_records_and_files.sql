-- Create records table for activity logging
CREATE TABLE records (
    id BIGSERIAL PRIMARY KEY,
    object_id BIGINT NOT NULL REFERENCES object_definitions(id) ON DELETE RESTRICT,
    data JSONB NOT NULL,
    created_by_id BIGINT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create files table for photo/video uploads
CREATE TABLE files (
    id BIGSERIAL PRIMARY KEY,
    record_id BIGINT NOT NULL REFERENCES records(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(1024) NOT NULL,
    file_type VARCHAR(50),
    file_size BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_records_object_id ON records(object_id);
CREATE INDEX idx_records_created_by_id ON records(created_by_id);
CREATE INDEX idx_records_created_at ON records(created_at DESC);
CREATE INDEX idx_records_deleted_at ON records(deleted_at);
CREATE INDEX idx_files_record_id ON files(record_id);
