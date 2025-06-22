CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL PRIMARY KEY,
    name text NOT NULL,
    email text NOT NULL
);