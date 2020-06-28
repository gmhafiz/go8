-- +goose Up
-- SQL in this section is executed when the migrations is applied.
CREATE table authors
(
    author_id varchar(20) not null,
    first_name varchar(255) not null,
    middle_name varchar(255),
    last_name varchar(255) not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    primary key (author_id)
);
-- +goose Down
-- SQL in this section is executed when the migrations is rolled back.
drop table authors;