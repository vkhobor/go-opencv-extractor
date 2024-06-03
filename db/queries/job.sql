-- name: ListJobsWithProgress :many
SELECT
    *
FROM
    jobs
    LEFT JOIN yt_videos ON jobs.id = yt_videos.job_id;

-- name: GetJob :one
SELECT
    j.id AS id,
    j.search_query AS search_query,
    j."limit" AS "limit",
    COUNT(DISTINCT v.id) AS videos_found
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
    jobs (id, search_query, "limit")
VALUES
    (?, ?, ?) RETURNING *;

-- name: GetJobs :many
SELECT
    COALESCE(found_videos, 0),
    jobs.id,
    jobs."limit",
    jobs.search_query,
    jobs.filter_id
FROM
    jobs
    LEFT JOIN (
        SELECT
            COUNT(*) as found_videos,
            jobs.id
        FROM
            jobs
            LEFT JOIN yt_videos ON jobs.id = yt_videos.job_id
        WHERE
            yt_videos.id IS NOT NULL
        GROUP BY
            jobs.id
    ) as t ON t.id = jobs.id;

-- name: GetOneWithVideos :many
SELECT
    j.id AS id,
    v.id AS video_youtube_id,
    COUNT(DISTINCT p.id) AS pictures_found
FROM
    jobs j
    LEFT JOIN yt_videos v ON j.id = v.job_id
    LEFT JOIN pictures p ON v.id = p.yt_video_id
WHERE
    j.id = ?
GROUP BY
    j.id,
    v.id;