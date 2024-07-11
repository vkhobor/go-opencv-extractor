import { Injectable, inject } from '@angular/core';
import { injectQuery, injectMutation, injectQueryClient } from '@ngneat/query';
import { Observable } from 'rxjs';
import { Job } from '../models/Job';
import { DefaultHttpProxyService } from './http/default-http-proxy.service';
import { Stats } from '../models/Stats';

@Injectable({
    providedIn: 'root',
})
export class StatsService {
    #http = inject(DefaultHttpProxyService);
    #query = injectQuery();
    #mutate = injectMutation();
    #queryClient = injectQueryClient();

    constructor() {}

    getStats() {
        return this.#query({
            queryKey: ['stats'] as const,
            refetchInterval: 5000,
            initialData: {
                macthing_pictures_saved: 0,
                videos_checked: 0,
            } as Stats,

            queryFn: () => {
                return this.#http.get('/stats') as Observable<Stats>;
            },
        });
    }
}
