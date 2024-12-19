CREATE TABLE urls(
                     user_id int,
                     short_url text,
                     long_url text,
                     PRIMARY KEY (user_id, short_url)
);
