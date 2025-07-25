-- name: AddYtVideo :one
INSERT INTO
    yt_videos (id, job_id, name)
VALUES
    (?, ?, ?) RETURNING *;

-- name: GetYtVideoWithJob :one
SELECT
    *
FROM
    yt_videos
    JOIN jobs ON yt_videos.job_id = jobs.id
WHERE
    yt_videos.id = ?;

-- name: GetVideosDownloadedButNotImported :many
SELECT
    yt_videos.id as yt_video_id,
    jobs.id AS job_id,
    jobs.search_query,
    jobs.filter_id,
    jobs."limit",
    jobs.youtube_id,
    blob_storage.path AS path,
    import_attempts.progress AS import_progress
FROM
    yt_videos
    JOIN download_attempts ON yt_videos.id = download_attempts.yt_video_id
    JOIN blob_storage ON download_attempts.blob_storage_id = blob_storage.id
    LEFT JOIN import_attempts ON yt_videos.id = import_attempts.yt_video_id
    JOIN jobs ON jobs.id = yt_videos.job_id
WHERE
    import_attempts.progress is not 100
GROUP BY
    yt_videos.id;

-- name: GetVideosDownloaded :many
SELECT
    yt_videos.id as yt_video_id,
    jobs.id AS job_id,
    jobs.search_query,
    jobs.filter_id,
    jobs."limit",
    jobs.youtube_id,
    yt_videos.name as yt_video_name,
    blob_storage.path AS path,
    CAST(
        COALESCE(MAX(import_attempts.progress), 0) AS INTEGER
    ) AS import_progress
FROM
    yt_videos
    JOIN download_attempts ON yt_videos.id = download_attempts.yt_video_id
    JOIN blob_storage ON download_attempts.blob_storage_id = blob_storage.id
    LEFT JOIN import_attempts ON yt_videos.id = import_attempts.yt_video_id
    JOIN jobs ON jobs.id = yt_videos.job_id
GROUP BY
    yt_videos.id,
    jobs.id,
    yt_videos.name,
    blob_storage.path;

-- name: GetVideoWithImportAttempts :many
SELECT
    *
FROM
    yt_videos
    JOIN import_attempts ON yt_videos.id = import_attempts.yt_video_id
WHERE
    yt_videos.id = ?;

-- name: GetVideoWithDownloadAttempts :many
SELECT
    *
FROM
    yt_videos
    JOIN download_attempts ON yt_videos.id = download_attempts.yt_video_id
WHERE
    yt_videos.id = ?;
