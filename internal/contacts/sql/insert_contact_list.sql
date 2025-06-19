insert into contact_lists (name, user_id)
values ($1, $2)
returning id, uuid, name, created_at, updated_at, user_id;