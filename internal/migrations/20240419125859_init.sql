-- +goose Up
-- +goose StatementBegin
CREATE TABLE persons (
    id BIGSERIAL PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    patronymic TEXT,
    sex TEXT,
    country TEXT,
    created_at timestamp,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE person;
-- +goose StatementEnd
