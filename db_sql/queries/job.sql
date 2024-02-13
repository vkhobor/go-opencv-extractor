
-- name: ListJobsWithProgress :many
SELECT * FROM jobs
LEFT JOIN yt_videos ON jobs.id = yt_videos.job_id;

-- name: CreateJob :one
INSERT INTO jobs (
  id, search_query, "limit"
) VALUES (
  ?, ?, ?
)
RETURNING *;

