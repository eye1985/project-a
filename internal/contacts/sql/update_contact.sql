update contact
set has_accepted = $1,
    accepted_at  = $2
where uuid = $3
  and user_id = $4;