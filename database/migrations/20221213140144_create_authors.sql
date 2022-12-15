-- +goose Up
-- +goose StatementBegin
create table IF not exists authors
(
    id bigserial
        constraint authors_pk
            primary key,
    first_name text not null,
    middle_name text,
    last_name text not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone
);


CREATE TRIGGER update_author_updated_at BEFORE UPDATE
    ON authors FOR EACH ROW EXECUTE PROCEDURE
    update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table authors
-- +goose StatementEnd
