update magic_links
set is_activated = true
where public_id = $1
  AND expires_at >= $2
  AND is_activated = false
returning email