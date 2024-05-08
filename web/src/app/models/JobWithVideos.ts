export interface Video {
  youtube_id: string;
  status: string;
  pictures_found: number;
}

export interface JobWithVideos {
  id: string;
  videos: Video[];
}
