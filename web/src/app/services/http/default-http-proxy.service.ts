import { Injectable } from '@angular/core';
import { HttpProxyService } from './http-proxy.service';
import { HttpClient } from '@angular/common/http';
import env from '../../../enviroments/enviroment';

@Injectable({
    providedIn: 'root',
})
export class DefaultHttpProxyService extends HttpProxyService {
    override getBaseUrl(): string {
        return env.api;
    }

    constructor(http: HttpClient) {
        super(http);
    }
}
