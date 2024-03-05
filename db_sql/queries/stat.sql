-- name: ListImportedVideosWithSaved :many
SELECT yt_videos.id, COUNT(pictures.id) as 'importedPictures' FROM yt_videos
LEFT JOIN pictures ON yt_videos.id = pictures.yt_video_id
GROUP BY yt_videos.id
HAVING yt_videos.status = 'imported';