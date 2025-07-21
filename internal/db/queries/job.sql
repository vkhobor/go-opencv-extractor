-- name: GetJob :one
SELECT
    j.id AS id,
    j.search_query AS search_query,
    j."limit" AS "limit",
    j.youtube_id AS youtube_id,
    COUNT(v.id) AS videos_found
FROM
    jobs j
    LEFT JOIN yt_videos v ON j.id = v.job_id
WHERE
    j.id = ?
GROUP BY
    j.id,
    j.search_query,
    j."limit";

-- name: CreateJob :one
INSERT INTO
    jobs (id, search_query, "limit",youtube_id, filter_id)
VALUES
    (?, ?, ?,?, ?) RETURNING *;
