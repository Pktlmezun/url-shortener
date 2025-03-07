CREATE TABLE urls (
    user_id INT NOT NULL,
    short_url TEXT NOT NULL,
    long_url TEXT NOT NULL,
    PRIMARY KEY (user_id, short_url)
);
