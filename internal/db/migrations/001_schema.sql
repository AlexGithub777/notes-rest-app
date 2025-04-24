-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    category INTEGER REFERENCES categories(id),
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- Insert default user
INSERT INTO users (username, password)
VALUES ('user1', 'user1pass');  -- Store plain text passwords only for now, will hash it later!

-- Insert default categories
INSERT INTO categories (name) VALUES ('General'), ('Work'), ('Personal'), ('Other');

INSERT INTO notes (title, content, category, user_id, created_at, updated_at)
VALUES 
    ('First Note', 'This is the content of the first note', 1, 1, NOW(), NOW()),
    ('Second Note', 'This is the content of the second note', 2, 1, NOW(), NOW());

-- +goose Down
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
