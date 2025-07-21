-- name: AddBlob :exec
INSERT into
    blob_storage (id, path)
VALUES
    (?, ?);


-- name: GetBlob :one
SELECT
    path
from
    blob_storage
WHERE
    id = ?;
