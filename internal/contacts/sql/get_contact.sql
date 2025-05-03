SELECT distinct u.id, u.uuid as user_uuid, u.username, u.email, cl.name as list_name
FROM contact_list_link cll
         JOIN contact_lists cl ON cll.contact_list_id = cl.id
         JOIN contact c ON cll.contact_id = c.id
         JOIN users u ON u.id =
                         CASE
                             WHEN c.user_1 = cl.user_id THEN c.user_2
                             ELSE c.user_1
                             END
WHERE cl.id = $1
  AND c.removed_at IS NULL
ORDER BY u.username