-- +goose Up
-- +goose StatementBegin
CREATE TABLE userss (
    id SERIAL PRIMARY KEY,
    name TEXT,
    username TEXT UNIQUE NOT NULL,
    password TEXT UNIQUE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users
-- +goose StatementEnd
