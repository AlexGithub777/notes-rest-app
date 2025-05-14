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
VALUES 
    ('user1', 'user1pass'),
    ('user2', 'user2pass'),
    ('user3', 'user3pass'),
    ('user4', 'user4pass'),
    ('user5', 'user5pass'),
    ('user6', 'user6pass'),
    ('user7', 'user7pass'),
    ('user8', 'user8pass'),
    ('user9', 'user9pass'),
    ('user10', 'user10pass');

-- Insert default categories
INSERT INTO categories (name) VALUES ('General'), ('Work'), ('Personal'), ('Other');

INSERT INTO notes (title, content, category, user_id, created_at, updated_at)
VALUES 
    ('First Note', 'This is the content of the first note', 1, 1, NOW(), NOW()),
    ('Second Note', 'This is the content of the second note', 2, 2, NOW(), NOW()),
    ('Third Note', 'This is the content of the third note', 3, 3, NOW(), NOW()),
    ('Fourth Note', 'This is the content of the fourth note', 4, 4, NOW(), NOW()),
    ('Fifth Note', 'This is the content of the fifth note', 1, 5, NOW(), NOW()),
    ('Sixth Note', 'This is the content of the sixth note', 2, 6, NOW(), NOW()),
    ('Seventh Note', 'This is the content of the seventh note', 3, 7, NOW(), NOW()),
    ('Eighth Note', 'This is the content of the eighth note', 4, 8, NOW(), NOW()),
    ('Ninth Note', 'This is the content of the ninth note', 1, 9, NOW(), NOW()),
    ('Tenth Note', 'This is the content of the tenth note', 2, 10, NOW(), NOW());

-- +goose Down
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;
