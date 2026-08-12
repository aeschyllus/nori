-- +goose Up
CREATE TABLE IF NOT EXISTS categories (
  category_id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL
);

-- +goose Down
DROP TABLE categories;
