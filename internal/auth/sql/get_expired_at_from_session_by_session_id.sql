SELECT expires_at
FROM user_sessions
WHERE session_id = $1