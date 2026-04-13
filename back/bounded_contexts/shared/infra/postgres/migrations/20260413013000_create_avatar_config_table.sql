-- +goose Up
CREATE TABLE avatar_configs (
    avatar_id TEXT PRIMARY KEY REFERENCES avatars(id) ON DELETE CASCADE,
    prompt TEXT NOT NULL,
    artistic_style TEXT NOT NULL CHECK (artistic_style IN ('2D', '3D')),
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- +goose Down
DROP TABLE avatar_configs;
