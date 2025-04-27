select id, name, created_at, updated_at, user_id
from user_lists
where id = $1;