-- +goose Up
-- +goose StatementBegin
CREATE TABLE baseball (
    id SERIAL PRIMARY KEY,
    sport TEXT NOT NULL,
    rate  DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL
);
CREATE TABLE football (
    id SERIAL PRIMARY KEY,
    sport TEXT NOT NULL,
    rate  DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL
);
CREATE TABLE soccer (
    id SERIAL PRIMARY KEY,
    sport TEXT NOT NULL,
    rate  DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE baseball;
DROP TABLE football;
DROP TABLE soccer;
-- +goose StatementEnd
