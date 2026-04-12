-- +goose Up
CREATE TABLE avatars (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id),
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX idx_avatars_user_id ON avatars (user_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_avatars_user_created_at ON avatars (user_id, created_at DESC) WHERE deleted_at IS NULL;

-- +goose Down
DROP TABLE avatars;
