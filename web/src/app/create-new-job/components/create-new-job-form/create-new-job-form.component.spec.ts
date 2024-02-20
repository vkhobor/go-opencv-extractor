import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateNewJobFormComponent } from './create-new-job-form.component';

describe('CreateNewJobFormComponent', () => {
  let component: CreateNewJobFormComponent;
  let fixture: ComponentFixture<CreateNewJobFormComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateNewJobFormComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(CreateNewJobFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
