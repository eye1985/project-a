select id, name, created_at, updated_at, user_id
from contact_lists
where user_id = $1;