import { Injectable } from '@angular/core';
import {
    InfiniteData,
    InfiniteQueryObserverResult,
    injectInfiniteQuery,
    injectQuery,
    QueryKey,
    QueryObserverResult,
} from '@ngneat/query';
import { client } from './http/kiota';

import { undefToErr } from './http/undefToErr';
import { CreateInfiniteQueryOptions } from '@ngneat/query/lib/infinite-query';
import { Response } from '../../api/models';
import { Result } from '@ngneat/query/lib/types';
import { ImagesRequestBuilderGetQueryParameters } from '../../api/api/images';

export type ImagePageQueryParams = Omit<
    ImagesRequestBuilderGetQueryParameters,
    'limit' | 'offset'
>;

@Injectable({
    providedIn: 'root',
})
export class ImagesService {
    #query = injectQuery();

    getImagePage(
        pageParam: number,
        pageSize: number,
        params: ImagePageQueryParams
    ) {
        return this.#query({
            queryKey: ['images', params, pageParam, pageSize] as const,
            enabled: false,
            queryFn: () => this.getImagePageApi(pageParam, pageSize, params),
        });
    }

    getImagePageApi(
        pageNumber: number,
        pageSize: number,
        params: ImagePageQueryParams
    ) {
        const offset = pageNumber * pageSize;
        const limit = pageSize;

        return undefToErr(
            client.api.images.get({
                queryParameters: {
                    limit,
                    offset,
                    ...params,
                },
            })
        );
    }
}
