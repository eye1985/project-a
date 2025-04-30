update contact_lists
set name=$1,
    updated_at=$2
where id = $3;