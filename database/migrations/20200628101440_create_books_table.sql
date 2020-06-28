-- +goose Up
-- SQL in this section is executed when the migrations is applied.
CREATE table books
(
    book_id varchar(20) not null,
    title varchar(255) not null,
    published_date timestamp with time zone not null,
    image_url varchar(255),
    description text,
    created_at timestamp with time zone default current_timestamp,
	updated_at timestamp with time zone default current_timestamp,
	deleted_at timestamp with time zone ,
	primary key (book_id)
);

-- +goose Down
-- SQL in this section is executed when the migrations is rolled back.
drop table books;