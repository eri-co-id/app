-- Create object_definitions table (schema-driven architecture)
CREATE TABLE object_definitions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    label VARCHAR(255) NOT NULL,
    plural_label VARCHAR(255) NOT NULL,
    icon VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create field_definitions table
CREATE TABLE field_definitions (
    id BIGSERIAL PRIMARY KEY,
    object_id BIGINT NOT NULL REFERENCES object_definitions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    label VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    required BOOLEAN DEFAULT FALSE,
    position INT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(object_id, name)
);

CREATE INDEX idx_object_definitions_name ON object_definitions(name);
CREATE INDEX idx_field_definitions_object_id ON field_definitions(object_id);
CREATE INDEX idx_field_definitions_type ON field_definitions(type);
