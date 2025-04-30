SELECT c.id,
       c.uuid,
       c.invited_by,
       inviter.email    AS invited_by_email,
       inviter.username AS inviter_username,
       c.user_id,
       invitee.email    AS invitee_email,
       invitee.username AS invitee_username,
       c.has_accepted
FROM contact c
         JOIN users invitee ON c.user_id = invitee.id
         JOIN users inviter ON c.invited_by = inviter.id
WHERE ($1 IN (c.user_id, c.invited_by));