create table IF not exists authors
(
    id bigserial
        constraint authors_pk
            primary key,
    first_name text not null,
    middle_name text,
    last_name text not null,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone
);

alter table authors owner to "user";

