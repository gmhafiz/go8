create table if not exists book_authors
(
    book_id bigserial not null,
    author_id bigserial not null,
    constraint book_authors_pk
        primary key (book_id, author_id)
);

