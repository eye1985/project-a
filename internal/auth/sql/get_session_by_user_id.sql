SELECT user_id, session_id, expires_at
FROM user_sessions
WHERE user_id = $1