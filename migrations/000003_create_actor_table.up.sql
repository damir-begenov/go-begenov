CREATE TABLE IF NOT EXISTS actor (
    -- id column is a 64-bit auto-incrementing integer & primary key (defines the row)
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone not null default NOW(),
    fullname text not null,
    year integer not null,
    films text[] not NULL,
    girlfriend text
    );