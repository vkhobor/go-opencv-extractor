import { Injectable, inject } from '@angular/core';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
  injectInfiniteQuery,
} from '@ngneat/query';
import { Observable, map } from 'rxjs';
import { DefaultHttpProxyService } from './http/default-http-proxy.service';
import { ImagesResponse } from '../models/Image';

@Injectable({
  providedIn: 'root',
})
export class ImagesService {
  #http = inject(DefaultHttpProxyService);
  #query = injectQuery();
  #mutate = injectMutation();
  #queryClient = injectQueryClient();
  #infinite = injectInfiniteQuery();

  getImages(pageSize: number) {
    return this.#infinite({
      queryKey: ['images'] as const,
      queryFn: ({ pageParam }) => this.getImagePage(pageParam, pageSize),
      initialPageParam: 0,
      getPreviousPageParam: (firstPage, _, params) => {
        if (params === 0) return undefined;
        return params - 1;
      },
      getNextPageParam: (lastPage, allPages, params) => {
        if (allPages.map((x) => x.pictures).flat().length >= lastPage.total)
        return undefined;
        return params + 1;
      },
    });
  }

  getImagePage(pageNumber: number, pageSize: number) {
    const offset = pageNumber * pageSize;
    const limit = pageSize;

    return this.#http.get(
      `/images?offset=${offset}&limit=${limit}`
    ) as Observable<ImagesResponse>;
  }
}
