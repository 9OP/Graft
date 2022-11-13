CREATE TABLE IF NOT EXISTS users (
    name TEXT NOT NULL UNIQUE
);

INSERT INTO users (name)
VALUES
    ("turing"),
    ("ritchie"),
    ("lovelace"),
    ("hamilton");
