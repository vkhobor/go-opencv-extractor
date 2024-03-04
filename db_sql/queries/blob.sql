-- name: AddBlob :one
INSERT into blob_storage
    (id, path)
    VALUES (?, ?)
RETURNING *;

-- name: GetBlob :one
SELECT path from blob_storage
    WHERE id = ?;
     