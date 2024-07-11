import { TestBed } from '@angular/core/testing';

import { DefaultHttpProxyService } from './default-http-proxy.service';

describe('DefaultHttpProxyService', () => {
    let service: DefaultHttpProxyService;

    beforeEach(() => {
        TestBed.configureTestingModule({});
        service = TestBed.inject(DefaultHttpProxyService);
    });

    it('should be created', () => {
        expect(service).toBeTruthy();
    });
});
