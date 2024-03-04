-- name: AddReference :exec
INSERT into reference_images
(blob_storage_id)
VALUES (sqlc.arg(id))
RETURNING *;

-- name: GetReferences :many
SELECT * FROM reference_images JOIN blob_storage ON reference_images.blob_storage_id = blob_storage.id;

-- name: DeleteReferences :exec
DELETE FROM reference_images;




