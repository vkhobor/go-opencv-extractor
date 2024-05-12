export interface Filter {
  id: string;
  name: string;
  filter_images: FilterImage[];
}

export interface FilterImage {
  blob_id: string;
}
