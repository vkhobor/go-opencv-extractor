import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddSearchModalComponent } from './add-search-modal.component';

describe('AddSearchModalComponent', () => {
  let component: AddSearchModalComponent;
  let fixture: ComponentFixture<AddSearchModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddSearchModalComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(AddSearchModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
