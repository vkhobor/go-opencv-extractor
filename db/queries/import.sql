-- name: AddImportAttempt :one
INSERT INTO import_attempts (
  id, yt_video_id, filter_id, progress, error
) VALUES (
  ?, ?, ?, ?, ?
) RETURNING *;