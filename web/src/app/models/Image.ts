export interface Image {
    id: string;
    blob_id: string;
    youtube_id: string;
}

export interface ImagesResponse {
    pictures: Image[];
    total: number;
}
