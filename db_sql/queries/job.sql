
-- name: ListJobsWithProgress :many
SELECT * FROM jobs
LEFT JOIN yt_videos ON jobs.id = yt_videos.job_id
LEFT JOIN pictures ON yt_videos.id = pictures.yt_video_id;

-- name: CreateJob :one
INSERT INTO jobs (
  id, search_query, "limit"
) VALUES (
  ?, ?, ?
)
RETURNING *;

-- name: GetToScrapeVideos :many
SELECT COALESCE(found_videos, 0), jobs.id, jobs."limit", jobs.search_query  FROM jobs
LEFT JOIN
(SELECT COUNT(*) as found_videos, jobs.id FROM jobs
LEFT JOIN yt_videos ON jobs.id = yt_videos.job_id
WHERE yt_videos.id IS NOT NULL
GROUP BY jobs.id) as t ON t.id = jobs.id

