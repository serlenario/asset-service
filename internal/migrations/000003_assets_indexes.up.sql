-- +goose Up
CREATE INDEX idx_assets_uid_name ON assets(uid, name);
CREATE INDEX idx_assets_uid_created_at ON assets(uid, created_at DESC);