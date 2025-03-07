CREATE TABLE counter (
    id SERIAL PRIMARY KEY,
    counter BIGINT NOT NULL DEFAULT 0
);

-- Insert initial counter record
INSERT INTO
    counter (counter)
VALUES
    (0);
