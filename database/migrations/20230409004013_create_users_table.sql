-- +goose Up
-- +goose StatementBegin
CREATE TABLE users
(
    id          bigint generated always as identity primary key ,
    first_name  text,
    middle_name text,
    last_name   text,
    email       text unique,
    password    text,
    verified_at timestamptz
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
