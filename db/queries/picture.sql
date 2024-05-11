
-- name: AddPicture :one
INSERT into pictures
    (id, yt_video_id, frame_number, blob_storage_id)
    VALUES (?,?,?,?)
RETURNING *;


-- name: GetPictures :many
SELECT * FROM pictures 
LIMIT ? OFFSET ?;

-- name: AllPicturesCount :one
SELECT COUNT(*) as count_all FROM pictures;