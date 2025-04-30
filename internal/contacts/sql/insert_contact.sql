insert into contact (user_id, invited_by, display_name, list_id)
values ($1, $2, $3, $4)
returning
    id,
    user_id,
    invited_by,
    has_accepted,
        (select email from users where id = $1) as invitee_email,
        (select email from users where id = $2) as inviter_email;