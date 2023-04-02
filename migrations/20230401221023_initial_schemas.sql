-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pages (
    key TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    next_page_key TEXT REFERENCES pages(key) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS pages_created_at ON pages(created_at);

CREATE TABLE IF NOT EXISTS lists (
    key TEXT PRIMARY KEY,
    next_page_key TEXT REFERENCES pages(key) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE lists;

DROP TABLE pages;

-- +goose StatementEnd