create table book_authors
(
    books_id bigserial
        constraint books_authors_books_book_id_fk
            references books,
    author_id bigserial
        constraint books_authors_authors_author_id_fk
            references authors,
    constraint books_authors_pkey
        primary key (books_id, author_id)
);