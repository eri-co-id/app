-- Create equipment table (master data per SPEC.md §3.2/§5)
CREATE TABLE equipment (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    code VARCHAR(50) UNIQUE,
    serial_batch VARCHAR(100),
    condition VARCHAR(50) NOT NULL DEFAULT 'good'
        CHECK (condition IN ('new', 'good', 'worn', 'damaged')),
    current_horse_id BIGINT REFERENCES horses(id),
    location VARCHAR(255),
    purchase_date DATE,
    notes TEXT,
    photo_url VARCHAR(2048),
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_equipment_current_horse_id ON equipment(current_horse_id);
CREATE INDEX idx_equipment_active ON equipment(active);
CREATE INDEX idx_equipment_deleted_at ON equipment(deleted_at);
