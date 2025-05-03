select invites.id,
       invites.uuid,
       invites.inviter_id,
       invites.invitee_id,
       inviter.email as inviter_email,
       invitee.email as invitee_email,
       invites.has_accepted
from invites
         join users as inviter on inviter.id = invites.inviter_id
         join users as invitee on invitee.id = invites.invitee_id
where has_accepted = false
  and $1 IN (invites.invitee_id, invites.inviter_id);