-- name: AddDownloadAttempt :exec
INSERT INTO
  download_attempts (id, yt_video_id, error, blob_storage_id, progress)
VALUES
  (?, ?, ?, ?, ?);