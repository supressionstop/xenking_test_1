-- +goose Up
-- +goose StatementBegin
CREATE TABLE lines (
    id SERIAL PRIMARY KEY,
    sport TEXT NOT NULL,
    rate  DECIMAL NOT NULL,
    created_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE lines;
-- +goose StatementEnd
