-- name: AddPicture :one
INSERT into
    pictures (
        id,
        import_attempt_id,
        frame_number,
        blob_storage_id
    )
VALUES
    (?, ?, ?, ?) RETURNING *;

-- name: GetPictures :many
SELECT
    *
FROM
    pictures
    JOIN import_attempts ON pictures.import_attempt_id = import_attempts.id
LIMIT
    ?
OFFSET
    ?;

-- name: AllPicturesCount :one
SELECT
    COUNT(*) as count_all
FROM
    pictures;