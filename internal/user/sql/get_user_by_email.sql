select id, username, email, created_at
from users
where email = $1