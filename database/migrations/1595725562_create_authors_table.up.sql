CREATE table authors
(
    author_id bigserial,
    first_name varchar(255) not null,
    middle_name varchar(255),
    last_name varchar(255) not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    deleted_at timestamp with time zone,
    primary key (author_id)
);