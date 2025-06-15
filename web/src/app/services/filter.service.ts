import { Injectable, inject } from '@angular/core';
import { DefaultHttpProxyService } from './http/default-http-proxy.service';
import {
    injectQuery,
    injectMutation,
    injectQueryClient,
    injectInfiniteQuery,
    QueryObserverResult,
} from '@ngneat/query';
import { Filter } from '../models/Filter';
import { Observable, of } from 'rxjs';

@Injectable({
    providedIn: 'root',
})
export class FilterService {
    #http = inject(DefaultHttpProxyService);
    #query = injectQuery();

    getFilters() {
        return this.#query({
            queryKey: ['filters'] as const,
            queryFn: () => {
                return this.#http.get(`/filters`) as Observable<Filter[]>;
            },
        });
    }

    private getFilterRaw(id: string) {
        return this.#http.get(`/filters/${id}`) as Observable<Filter>;
    }

    getFilter(id?: string) {
        return this.#query({
            queryKey: ['filters', id],
            queryFn: (params) => this.getFilterRaw(params.queryKey[1]!),
            enabled: !!id,
        });
    }

    addFilter() {
        // add filter logic here
    }

    constructor() {}
}
