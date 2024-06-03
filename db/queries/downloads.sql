-- name: AddDownloadAttempt :one

INSERT INTO download_attempts (
  yt_video_id, error, blob_storage_id, progress
) VALUES (
  ?, ?,?,?
)
RETURNING *;
