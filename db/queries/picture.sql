-- name: AddPicture :exec
INSERT into
    pictures (
        id,
        import_attempt_id,
        frame_number,
        blob_storage_id
    )
VALUES
    (?, ?, ?, ?);

-- name: GetPictures :many
SELECT
    *
FROM
    pictures
    JOIN import_attempts ON pictures.import_attempt_id = import_attempts.id
WHERE
    @is_filter_by_youtube_id = false
    OR import_attempts.yt_video_id = @youtube_id
LIMIT
    ?
OFFSET
    ?;

-- name: AllPicturesCount :one
SELECT
    COUNT(*) as count_all
FROM
    pictures
    JOIN import_attempts ON pictures.import_attempt_id = import_attempts.id
WHERE
    @is_filter_by_youtube_id = false
    OR import_attempts.yt_video_id = @youtube_id;