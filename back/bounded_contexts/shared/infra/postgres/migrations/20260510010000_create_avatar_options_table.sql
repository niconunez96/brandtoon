-- +goose Up
CREATE TABLE avatar_options (
    id text PRIMARY KEY,
    avatar_id text NOT NULL REFERENCES avatars(id),
    status text NOT NULL CHECK (status IN ('PENDING', 'DONE', 'FAILED')),
    href text NULL,
    selected boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    deleted_at timestamptz NULL
);

CREATE UNIQUE INDEX avatar_options_single_selected_idx
    ON avatar_options (avatar_id)
    WHERE selected = true AND deleted_at IS NULL;

CREATE INDEX avatar_options_avatar_id_created_at_idx
    ON avatar_options (avatar_id, created_at DESC)
    WHERE deleted_at IS NULL;

-- +goose Down
DROP INDEX IF EXISTS avatar_options_avatar_id_created_at_idx;
DROP INDEX IF EXISTS avatar_options_single_selected_idx;
DROP TABLE IF EXISTS avatar_options;
