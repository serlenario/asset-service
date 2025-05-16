-- +goose Down
DROP INDEX IF EXISTS idx_assets_uid_created_at;
DROP INDEX IF EXISTS idx_assets_uid_name;