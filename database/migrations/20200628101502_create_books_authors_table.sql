-- +goose Up
-- SQL in this section is executed when the migrations is applied.
CREATE table books_authors
(
    books_id varchar(20) not null,
    author_id varchar(20) not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone ,
    primary key (books_id, author_id)
);
-- +goose Down
-- SQL in this section is executed when the migrations is rolled back.
drop table books_authors;