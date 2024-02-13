-- name: AddBlob :one
INSERT into blob_storage
    (id, path)
    VALUES (?, ?)
RETURNING *;
