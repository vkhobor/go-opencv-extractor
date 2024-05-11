
-- name: ListJobsWithProgress :many
SELECT * FROM jobs
LEFT JOIN yt_videos ON jobs.id = yt_videos.job_id
LEFT JOIN pictures ON yt_videos.id = pictures.yt_video_id;

-- name: GetJob :one
SELECT
    j.id AS id,
    j.search_query AS search_query,
    j."limit" AS "limit",
    COUNT(DISTINCT v.id) AS videos_found,
    COUNT(DISTINCT CASE WHEN v.status = 'errored' THEN v.id END) AS videos_in_error,
    COUNT(DISTINCT p.id) AS pictures_found
FROM
    jobs j
LEFT JOIN
    yt_videos v ON j.id = v.job_id
LEFT JOIN
    pictures p ON v.id = p.yt_video_id
WHERE
    j.id = ?
GROUP BY
    j.id, j.search_query, j."limit";

-- name: GetJobWithProgress :one
SELECT
    j.id AS id,
    j."limit" AS "limit",
    COUNT(DISTINCT v.id) AS videos_found,
    COUNT(DISTINCT CASE WHEN v.status = 'errored' THEN v.id END) AS videos_downloaded,
    COUNT(DISTINCT CASE WHEN v.status = 'imported' THEN v.id END) AS videos_imported,
    COUNT(DISTINCT CASE WHEN v.status = 'scraped' THEN v.id END) AS videos_scraped,
    COUNT(DISTINCT p.id) AS pictures_found
FROM
    jobs j
LEFT JOIN
    yt_videos v ON j.id = v.job_id
LEFT JOIN
    pictures p ON v.id = p.yt_video_id
WHERE
    j.id = ?
GROUP BY
    j.id, j.search_query, j."limit";

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
GROUP BY jobs.id) as t ON t.id = jobs.id;

-- name: GetOneWithVideos :many
SELECT
    j.id AS id,
    v.id AS video_youtube_id,
    v.status AS video_status,
    COUNT(DISTINCT p.id) AS pictures_found
FROM
    jobs j
LEFT JOIN
    yt_videos v ON j.id = v.job_id
LEFT JOIN
    pictures p ON v.id = p.yt_video_id
WHERE
    j.id = ?
GROUP BY
    j.id, v.id;
