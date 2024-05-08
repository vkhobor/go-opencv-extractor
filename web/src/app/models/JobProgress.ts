export interface Progress {
  imported: number;
  scraped: number;
  downloaded: number;
  video_ids: string[];
  number_of_pictures: number;
  limit: number;
}
