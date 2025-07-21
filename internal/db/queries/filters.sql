-- name: AddFilter :one
INSERT
OR REPLACE INTO filters (
    id,
    "name",
    discriminator,
    ratioTestThreshold,
    minThresholdForSURFMatches,
    minSURFMatches,
    MSESkip
)
VALUES
    (?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: AttachImageToFilter :one
INSERT INTO
    filter_images (filter_id, blob_storage_id)
VALUES
    (?, ?) RETURNING *;

-- name: DeleteImagesOnFilter :exec
DELETE FROM filter_images
WHERE
    filter_id = sqlc.arg (filter_id);

DELETE FROM blob_storage
WHERE
    blob_storage.id IN (
        SELECT
            blob_storage2.id
        FROM
            blob_storage blob_storage2
            JOIN filter_images ON filter_images.blob_storage_id = blob_storage2.id
        WHERE
            filter_images.filter_id = sqlc.arg (filter_id)
    );

-- name: GetFilters :many
SELECT
    blob_storage.id as blob_id,
    *
FROM
    filters
    LEFT JOIN filter_images ON filters.id = filter_images.filter_id
    LEFT JOIN blob_storage ON filter_images.blob_storage_id = blob_storage.id;

-- name: GetFilterById :many
SELECT
    f.*,
    fi.blob_storage_id
FROM
    filters f
    LEFT JOIN filter_images fi ON f.id = fi.filter_id
WHERE
    f.id = ?;

-- name: GetFilterForJob :many
SELECT
    filters.*,
    blob_storage.*
FROM
    jobs
    JOIN filters ON jobs.filter_id = filters.id
    JOIN filter_images ON filters.id = filter_images.filter_id
    JOIN blob_storage ON filter_images.blob_storage_id = blob_storage.id
WHERE
    jobs.id = ?;
