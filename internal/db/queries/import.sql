-- name: AddImportAttempt :one
INSERT INTO
  import_attempts (id, yt_video_id, filter_id, progress, error)
VALUES
  (?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateImportAttemptProgress :exec
UPDATE import_attempts
SET
  progress = ?
WHERE
  id = ?;

-- name: UpdateImportAttemptError :exec
UPDATE import_attempts
SET
  error = ?
WHERE
  id = ?;
