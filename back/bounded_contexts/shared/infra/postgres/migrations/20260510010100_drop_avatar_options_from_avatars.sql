-- +goose Up
ALTER TABLE avatars
    DROP COLUMN IF EXISTS avatar_options;

-- +goose Down
ALTER TABLE avatars
    ADD COLUMN avatar_options jsonb[] NOT NULL DEFAULT ARRAY[]::jsonb[];
