CREATE TABLE IF NOT EXISTS books
(
    book_id bigserial,
    title varchar(255) not null,
    published_date timestamp with time zone not null,
    image_url varchar(255),
    description text not null,
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    primary key (book_id)
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_book_updated_at BEFORE UPDATE
    ON books FOR EACH ROW EXECUTE PROCEDURE
    update_updated_at_column();
