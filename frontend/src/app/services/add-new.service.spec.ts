import { TestBed } from '@angular/core/testing';

import { AddNewService } from './add-new.service';

describe('AddNewService', () => {
  let service: AddNewService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(AddNewService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
