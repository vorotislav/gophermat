CREATE TABLE users (
    id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
    login        TEXT NOT NULL UNIQUE,
    password     TEXT NOT NULL
);

CREATE INDEX user_login_idx ON users (login);