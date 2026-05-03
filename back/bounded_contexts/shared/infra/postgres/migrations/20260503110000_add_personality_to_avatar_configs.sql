-- +goose Up
ALTER TABLE avatar_configs
    ADD COLUMN personality TEXT NOT NULL DEFAULT 'Friendly';

ALTER TABLE avatar_configs
    ADD CONSTRAINT avatar_configs_personality_check
    CHECK (personality IN ('Friendly', 'Bold', 'Playful'));

ALTER TABLE avatar_configs
    ALTER COLUMN personality DROP DEFAULT;

-- +goose Down
ALTER TABLE avatar_configs
    DROP CONSTRAINT IF EXISTS avatar_configs_personality_check;

ALTER TABLE avatar_configs
    DROP COLUMN IF EXISTS personality;
