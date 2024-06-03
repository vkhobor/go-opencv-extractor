-- name: AddYtVideo :one
INSERT INTO
  yt_videos (id, job_id)
VALUES
  (?, ?) RETURNING *;

-- name: GetYtVideo :one
SELECT
  *
FROM
  yt_videos
WHERE
  id = ?;

-- name: GetScrapedVideos :many
SELECT
  *
FROM
  yt_videos
WHERE
  status = "scraped";

-- name: GetVideosDownloaded :many
SELECT
  *
FROM
  yt_videos
  JOIN blob_storage ON yt_videos.blob_storage_id = blob_storage.id
WHERE
  status = "downloaded";

-- name: GetJobVideosWithProgress :many
SELECT
  *
FROM
  yt_videos
  LEFT JOIN download_attempts ON yt_videos.id = download_attempts.yt_video_id
  LEFT JOIN import_attempts ON yt_videos.id = import_attempts.yt_video_id
WHERE
  yt_videos.job_id = ?;