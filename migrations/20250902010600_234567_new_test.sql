-- +goose up

INSERT INTO todos (title, description, completed, due_date)
VALUES ('Test Todo from Goose', 'This is a test insert using Goose migration', FALSE, '2025-09-15');
VALUES ('Test Todo from Goose3', 'This is a test insert using Goose migration', FALSE, '2025-09-15');
VALUES ('Test Todo from Goose4', 'This is a test insert using Goose migration', FALSE, '2025-09-15');

-- +goose down 
