-- +goose Up
CREATE TABLE sessions (
    id TEXT PRIMARY KEY DEFAULT encode(gen_random_bytes(16), 'hex'),
    uid BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ip_addr INET NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    UNIQUE(uid)
);