-- Track who uploaded each file (SPEC.md §5 example has uploaded_by_id; 004 didn't)
ALTER TABLE files ADD COLUMN uploaded_by_id BIGINT REFERENCES users(id);
