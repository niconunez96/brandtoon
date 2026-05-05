-- +goose Up
ALTER TABLE avatars
    ADD COLUMN avatar_options jsonb[] NOT NULL DEFAULT ARRAY[]::jsonb[];

-- +goose Down
ALTER TABLE avatars
    DROP COLUMN IF EXISTS avatar_options;
