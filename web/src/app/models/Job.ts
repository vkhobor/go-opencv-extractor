export interface Job {
  search_query: string;
  id: string;
  limit: number;
  progress: {
    imported: number;
    downloaded: number;
    scraped: number;
    video_ids: string[];
    number_of_pictures: number;
  };
}
