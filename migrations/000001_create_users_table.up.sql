CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    username   TEXT        NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE user_lists
(
    id         SERIAL PRIMARY KEY,
    name       TEXT    NOT NULL,
    created_at timestamptz DEFAULT now(),
    user_id    integer not null references users (id) on delete cascade,
    UNIQUE (user_id, name)
);

create table user_list_record
(
    id           SERIAL PRIMARY KEY,
    user_id      INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    invited_by   INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    display_name TEXT    NOT NULL,
    has_accepted bool        DEFAULT false,
    invited_at   timestamptz DEFAULT now(),
    accepted_at  timestamptz,
    removed_at   timestamptz,
    list_id      integer not null references user_lists (id) on DELETE CASCADE,
    UNIQUE (user_id, invited_by),
    CHECK (user_id IS DISTINCT FROM invited_by)
);

CREATE TABLE user_sessions
(
    id         SERIAL PRIMARY KEY,
    session_id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ NOT NULL,
    user_id    integer     NOT NULL REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE magic_links
(
    id           SERIAL PRIMARY KEY,
    created_at   TIMESTAMPTZ DEFAULT NOW(),
    expires_at   TIMESTAMPTZ NOT NULL,
    email        TEXT        NOT NULL,
    public_id    TEXT        NOT NULL UNIQUE,
    is_activated BOOL        DEFAULT false
)
