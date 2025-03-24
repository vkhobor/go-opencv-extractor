CREATE TABLE
    jobs (
        id TEXT PRIMARY KEY,
        search_query TEXT,
        filter_id TEXT,
        youtube_id TEXT,
        "limit" INT,
        FOREIGN KEY (filter_id) REFERENCES filters (id)
    );

CREATE TABLE
    yt_videos (
        id TEXT PRIMARY KEY,
        job_id TEXT,
        FOREIGN KEY (job_id) REFERENCES jobs (id)
    );

CREATE TABLE
    download_attempts (
        id TEXT PRIMARY KEY,
        yt_video_id TEXT,
        progress INT,
        blob_storage_id TEXT,
        error TEXT,
        FOREIGN KEY (blob_storage_id) REFERENCES blob_storage (id),
        FOREIGN KEY (yt_video_id) REFERENCES yt_videos (id)
    );

CREATE TABLE
    import_attempts (
        id TEXT PRIMARY KEY,
        yt_video_id TEXT,
        filter_id TEXT,
        progress INT,
        error TEXT,
        FOREIGN KEY (filter_id) REFERENCES filters (id),
        FOREIGN KEY (yt_video_id) REFERENCES yt_videos (id)
    );

CREATE TABLE
    filters (id TEXT PRIMARY KEY, name TEXT);

CREATE TABLE
    filter_images (
        filter_id TEXT,
        blob_storage_id TEXT,
        FOREIGN KEY (filter_id) REFERENCES filters (id),
        FOREIGN KEY (blob_storage_id) REFERENCES blob_storage (id)
    );

CREATE TABLE
    pictures (
        id TEXT PRIMARY KEY,
        import_attempt_id TEXT,
        frame_number INT,
        blob_storage_id TEXT,
        FOREIGN KEY (import_attempt_id) REFERENCES import_attempts (id),
        FOREIGN KEY (blob_storage_id) REFERENCES blob_storage (id)
    );

CREATE TABLE
    blob_storage (id TEXT PRIMARY KEY, path TEXT NOT NULL);
