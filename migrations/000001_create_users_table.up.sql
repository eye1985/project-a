CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID        NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    username   TEXT        NOT NULL,
    email      TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ                 DEFAULT NOW(),
    updated_at timestamptz
);

CREATE TABLE contact_lists
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID    NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    name       TEXT    NOT NULL,
    created_at TIMESTAMPTZ             DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    user_id    INTEGER not null references users (id) on delete CASCADE,
    UNIQUE (user_id, name)
);

CREATE TABLE contact
(
    id         SERIAL PRIMARY KEY,
    uuid       UUID    NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    user_1     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    user_2     INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ             DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    removed_at TIMESTAMPTZ,
    UNIQUE (user_1, user_2),
    CHECK (user_1 < user_2)
);

CREATE TABLE contact_list_link
(
    id              SERIAL PRIMARY KEY,
    uuid            UUID    NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    contact_id      INTEGER NOT NULL REFERENCES contact (id) ON DELETE CASCADE,
    contact_list_id INTEGER NOT NULL REFERENCES contact_lists (id) ON DELETE CASCADE,
    created_at      TIMESTAMPTZ             DEFAULT NOW(),
    UNIQUE (contact_id, contact_list_id)
);

CREATE TABLE invites
(
    id           SERIAL PRIMARY KEY,
    uuid         UUID    NOT NULL UNIQUE DEFAULT gen_random_uuid(),
    inviter_id   INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    invitee_id   INTEGER NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    has_accepted BOOL                    DEFAULT false,
    invited_at   timestamptz             DEFAULT now(),
    accepted_at  timestamptz,
    UNIQUE (inviter_id, invitee_id),
    CHECK (inviter_id IS DISTINCT FROM invitee_id)
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
);

CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS
$$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_users_trigger
    BEFORE UPDATE
    on users
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER update_contact_lists_trigger
    BEFORE UPDATE
    on contact_lists
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER update_contact_trigger
    BEFORE UPDATE
    on contact
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
