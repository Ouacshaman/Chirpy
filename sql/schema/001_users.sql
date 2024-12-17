-- +goose Up
CREATE TABLE users (
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
email TEXT UNIQUE NOT NULL,
hashed_password TEXT DEFAULT 'unset' NOT NULL
);

CREATE TABLE chirps(
id UUID PRIMARY KEY,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
body TEXT NOT NULL,
user_id UUID NOT NULL
    REFERENCES users(id)
    ON DELETE CASCADE
);

-- +goose Down
DROP TABLE IF EXISTS chirps;
DROP TABLE IF EXISTS users;
