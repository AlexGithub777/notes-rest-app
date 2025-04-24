-- +goose Up
CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

INSERT INTO notes (title, content, created_at, updated_at)
VALUES 
    ('First Note', 'This is the content of the first note', NOW(), NOW()),
    ('Second Note', 'This is the content of the second note', NOW(), NOW());

-- +goose Down
DROP TABLE IF EXISTS notes;
