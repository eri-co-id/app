-- Create horses table (master data, fixed structure per SPEC.md §3.2/§5)
CREATE TABLE horses (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) UNIQUE,
    breed VARCHAR(100),
    dob DATE,
    color VARCHAR(100),
    owner_name VARCHAR(255),
    current_trainer_id BIGINT REFERENCES users(id),
    health_status VARCHAR(50) NOT NULL DEFAULT 'active'
        CHECK (health_status IN ('active', 'injured', 'retired')),
    notes TEXT,
    photo_url VARCHAR(2048),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_horses_current_trainer_id ON horses(current_trainer_id);
CREATE INDEX idx_horses_health_status ON horses(health_status);
CREATE INDEX idx_horses_deleted_at ON horses(deleted_at);
