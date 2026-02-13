-- +goose Up
CREATE TABLE IF NOT EXISTS url_table (
   id BIGSERIAL PRIMARY KEY,
   original_url TEXT NOT NULL UNIQUE,
   short_url VARCHAR(10) NOT NULL UNIQUE UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS url_table;