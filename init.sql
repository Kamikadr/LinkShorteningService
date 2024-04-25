CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email VARCHAR NOT NULL UNIQUE,
    password TEXT NOT NULL,
    refresh_token TEXT
);
CREATE TABLE IF NOT EXISTS url(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    old_url TEXT NOT NULL UNIQUE,
    new_url TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX idx_user_url ON url(user_id, old_url)