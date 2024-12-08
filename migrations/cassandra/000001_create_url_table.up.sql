CREATE TABLE urls(
    id UUID PRIMARY KEY,
    user_id int,
    short_url text,
    long_url text
);
