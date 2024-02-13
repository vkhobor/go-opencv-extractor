CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    search_query TEXT,
    "limit" INT
);

CREATE TABLE yt_videos (
    id TEXT PRIMARY KEY,
    job_id TEXT,
    status TEXT,
    error TEXT,
    blob_storage_id TEXT,
    FOREIGN KEY (blob_storage_id) REFERENCES blob_storage(id)
    FOREIGN KEY (job_id) REFERENCES jobs(id)
);

CREATE TABLE pictures (
    id TEXT PRIMARY KEY,
    yt_video_id TEXT,
    frame_number INT,
    blob_storage_id TEXT,
    FOREIGN KEY (yt_video_id) REFERENCES jobs(id)
    FOREIGN KEY (blob_storage_id) REFERENCES blob_storage(id)
);

CREATE TABLE blob_storage (
    id TEXT PRIMARY KEY,
    path TEXT
);
