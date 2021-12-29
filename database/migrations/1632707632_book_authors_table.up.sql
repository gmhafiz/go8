create table if not exists book_authors
(
    book_id bigserial not null
        constraint book_authors_books_book_id_fk
            references books
            on delete cascade,
    author_id bigserial not null
        constraint book_authors_authors_id_fk
            references authors
            on delete cascade,
    constraint book_authors_pk
        primary key (book_id, author_id)
);
