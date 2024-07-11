import { HttpClient } from '@angular/common/http';

export abstract class HttpProxyService {
    constructor(private http: HttpClient) {}

    abstract getBaseUrl(): string;

    get = addBaseUrl(this.http.get.bind(this.http), this.getBaseUrl());
    post = addBaseUrl(this.http.post.bind(this.http), this.getBaseUrl());
    put = addBaseUrl(this.http.put.bind(this.http), this.getBaseUrl());
    delete = addBaseUrl(this.http.delete.bind(this.http), this.getBaseUrl());
    patch = addBaseUrl(this.http.patch.bind(this.http), this.getBaseUrl());
}

const addBaseUrl = <T extends Array<any>, U>(
    fn: (url: string, ...args: T) => U,
    baseUrl: string
) => {
    return (url: string, ...args: T): U => fn(baseUrl + url, ...args);
};
