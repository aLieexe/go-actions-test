CREATE TABLE IF NOT EXISTS users(
    username text NOT NULL,
    hashed_password text NOT NULL,
    id text NOT NULL PRIMARY KEY,
    email text NOT NULL UNIQUE
);

