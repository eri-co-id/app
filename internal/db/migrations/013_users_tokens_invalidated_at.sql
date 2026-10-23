-- Enables real logout: tokens issued before this timestamp are rejected,
-- even though they're stateless JWTs that would otherwise still verify.
ALTER TABLE users ADD COLUMN tokens_invalidated_at TIMESTAMP;
