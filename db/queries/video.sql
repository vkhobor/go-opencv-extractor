
-- name: AddYtVideo :one
INSERT INTO yt_videos (
  id, job_id, status

) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: UpdateStatus :one
UPDATE yt_videos
SET
  status = ?
,error = ?
WHERE id = ?
RETURNING *;

-- name: AddBlobToVideo :one
UPDATE yt_videos
SET
 blob_storage_id = ?
WHERE id = ?
RETURNING *;

-- name: GetYtVideo :one
SELECT * FROM yt_videos WHERE id = ?;

-- name: GetScrapedVideos :many
SELECT * FROM yt_videos WHERE status = "scraped";

-- name: GetVideosDownloaded :many
SELECT * FROM yt_videos JOIN blob_storage ON yt_videos.blob_storage_id = blob_storage.id WHERE status = "downloaded";
