export interface Job {
  search_query: string;
  id: string;
  limit: number;
  progress: {
    imported: number;
    downloaded: number;
    scraped: number;
  };
}
