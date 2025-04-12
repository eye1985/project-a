CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    username   TEXT        NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ default NOW()
);

CREATE INDEX idx_email on users (email);

CREATE TABLE user_connections
(
    id         SERIAL PRIMARY KEY,
    user1_id   integer NOT NULL,
    user2_id   integer NOT NULL,
    created_at TIMESTAMPTZ default NOW(),
    CONSTRAINT fk_user1 FOREIGN KEY (user1_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT fk_user2 FOREIGN KEY (user2_id) REFERENCES users (id) ON DELETE CASCADE,
    CHECK (user1_id <> user2_id)
);

CREATE UNIQUE INDEX unique_user_connection
    ON user_connections (LEAST(user1_id, user2_id), GREATEST(user1_id, user2_id));

CREATE INDEX idx_user1 ON user_connections (user1_id);
CREATE INDEX idx_user2 ON user_connections (user2_id);

CREATE TABLE user_sessions
(
    id         SERIAL PRIMARY KEY,
    session_id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ default NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    user_id    integer     NOT NULL REFERENCES users (id) ON DELETE CASCADE
);
