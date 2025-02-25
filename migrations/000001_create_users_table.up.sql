CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    username   TEXT NOT NULL,
    email      TEXT,
    created_at TIMESTAMPTZ default NOW()
)