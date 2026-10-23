-- Create otp_sessions table for OTP authentication
CREATE TABLE otp_sessions (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    otp_code VARCHAR(6) NOT NULL,
    attempt_count INT DEFAULT 0,
    max_attempts INT DEFAULT 5,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_otp_sessions_email ON otp_sessions(email);
CREATE INDEX idx_otp_sessions_expires_at ON otp_sessions(expires_at);
CREATE INDEX idx_otp_sessions_created_at ON otp_sessions(created_at DESC);
