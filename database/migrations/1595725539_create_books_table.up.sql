CREATE table books
(
    book_id bigserial,
    title varchar(255) not null,
    published_date timestamp with time zone not null,
    image_url varchar(255),
    description text not null ,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone ,
    primary key (book_id)
);