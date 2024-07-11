import { Injectable, inject } from '@angular/core';
import { injectQuery, injectMutation, injectQueryClient } from '@ngneat/query';
import { Observable, map } from 'rxjs';
import { Stats } from '../models/Stats';
import { DefaultHttpProxyService } from './http/default-http-proxy.service';
import { HttpResponse } from '@angular/common/http';
import { downloadBlob } from '../util/downloadFile';

@Injectable({
    providedIn: 'root',
})
export class ZipService {
    #http = inject(DefaultHttpProxyService);
    #query = injectQuery();
    #mutate = injectMutation();
    #queryClient = injectQueryClient();

    private getZipRaw() {
        return this.#http.get('/zipped', {
            responseType: 'blob' as 'json',
        }) as Observable<Blob>;
    }

    getZip() {
        return this.#query({
            queryKey: ['zip'] as const,
            queryFn: () => {
                return this.getZipRaw();
            },
        });
    }

    downloadZip() {
        this.getZipRaw().subscribe((zip) => {
            downloadBlob(zip, 'workspace.zip');
        });
    }
}
