import { Injectable } from '@angular/core';
import { HttpProxyService } from './http-proxy.service';
import { HttpClient } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class DefaultHttpProxyService extends HttpProxyService {
  override getBaseUrl(): string {
    return 'http://localhost:3010';
  }

  constructor(http: HttpClient) {
    super(http);
  }
}
