select u.id, u.uuid, u.username, u.email, u.created_at
from users u
         join user_sessions us on u.id = us.user_id
where session_id = $1