-- Create role_permissions table for access control
CREATE TABLE role_permissions (
    id BIGSERIAL PRIMARY KEY,
    role VARCHAR(50) NOT NULL,
    resource VARCHAR(255) NOT NULL,
    action VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(role, resource, action)
);

CREATE INDEX idx_role_permissions_role ON role_permissions(role);
CREATE INDEX idx_role_permissions_resource ON role_permissions(resource);

-- Create dashboard_definitions table for reporting
CREATE TABLE dashboard_definitions (
    id BIGSERIAL PRIMARY KEY,
    object_id BIGINT NOT NULL REFERENCES object_definitions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    label VARCHAR(255) NOT NULL,
    widgets JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    UNIQUE(object_id, name)
);

CREATE INDEX idx_dashboard_definitions_object_id ON dashboard_definitions(object_id);
