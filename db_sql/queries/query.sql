-- name: GetJob :one
SELECT * FROM jobs
WHERE search_query = ? LIMIT 1;

-- name: ListJobs :many
SELECT * FROM jobs
ORDER BY name;

-- name: CreateJob :one
INSERT INTO jobs (
  search_query
) VALUES (
  ?
)
RETURNING *;