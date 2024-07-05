-- name: ListJobsWithVideos :many
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
    jobs (id, search_query, "limit", filter_id)
VALUES
    (?, ?, ?, ?) RETURNING *;

-- name: GetJobs :many
SELECT
    COUNT(yt_videos.id) AS videos_found,
    jobs.id,
    jobs."limit",
    jobs.search_query,
    jobs.filter_id
FROM
    jobs
    LEFT JOIN yt_videos ON jobs.id = yt_videos.job_id
GROUP BY
    jobs.id;

-- name: GetVideosForJob :many
SELECT
    v.id AS video_youtube_id,
    COUNT(successDownload.yt_video_id) AS download_attempts_success,
    COUNT(errorDownload.yt_video_id) AS download_attempts_error,
    COUNT(successImport.yt_video_id) AS import_attempts_success,
    COUNT(successImport.yt_video_id) AS import_attempts_error
FROM
    jobs j
    JOIN yt_videos v ON j.id = v.job_id
    LEFT JOIN (
        SELECT
            yt_video_id
        FROM
            download_attempts
        WHERE
            error is null
    ) successDownload ON v.id = successDownload.yt_video_id
    LEFT JOIN (
        SELECT
            yt_video_id
        FROM
            download_attempts
        WHERE
            error is not null
    ) errorDownload ON v.id = errorDownload.yt_video_id
    LEFT JOIN (
        SELECT
            yt_video_id
        FROM
            import_attempts
        WHERE
            error is null
            and progress = 100
    ) successImport ON v.id = successImport.yt_video_id
    LEFT JOIN (
        SELECT
            yt_video_id
        FROM
            import_attempts
        WHERE
            error is not null
            and progress != 100
    ) errorImport ON v.id = successImport.yt_video_id
WHERE
    j.id = ?
GROUP BY
    v.id;