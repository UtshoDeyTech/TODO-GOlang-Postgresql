CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP
);

-- Add some sample data
INSERT INTO todos (title, description, completed, created_at, updated_at, completed_at)
VALUES
    ('Learn Go', 'Study Go programming language basics', true, NOW(), NOW(), NOW()),
    ('Build a REST API', 'Create a Todo API using Go and PostgreSQL', false, NOW(), NOW(), NULL),
    ('Learn Docker', 'Learn how to containerize applications', false, NOW(), NOW(), NULL);