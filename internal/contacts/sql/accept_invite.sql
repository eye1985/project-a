update invites
set has_accepted= true,
    accepted_at=now()
where uuid = $1
  AND invitee_id = $2
returning invites.id, inviter_id, invitee_id;