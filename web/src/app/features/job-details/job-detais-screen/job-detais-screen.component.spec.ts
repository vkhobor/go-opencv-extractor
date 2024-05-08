import { ComponentFixture, TestBed } from '@angular/core/testing';

import { JobDetaisScreenComponent } from './job-detais-screen.component';

describe('JobDetaisScreenComponent', () => {
  let component: JobDetaisScreenComponent;
  let fixture: ComponentFixture<JobDetaisScreenComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [JobDetaisScreenComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(JobDetaisScreenComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
