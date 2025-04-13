insert into users(username, email)
values ($1, $2)
returning id,email,username, created_at