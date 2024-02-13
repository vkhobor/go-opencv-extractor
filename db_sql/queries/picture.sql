
-- name: AddPicture :one
INSERT into pictures
    (id, yt_video_id, frame_number, blob_storage_id)
    VALUES (?,?,?,?)
RETURNING *;
