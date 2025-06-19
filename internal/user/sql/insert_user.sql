insert into users(username, email)
values ($1, $2)
returning id,uuid,email,username, created_at