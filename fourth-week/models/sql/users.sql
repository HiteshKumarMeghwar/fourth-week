CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT,
    username TEXT,
    password TEXT NOT NULL
);