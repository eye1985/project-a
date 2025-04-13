INSERT INTO user_sessions (user_id, session_id, expires_at)
VALUES ($1, $2, $3)
returning user_id, session_id, expires_at