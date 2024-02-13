import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProgressStateComponent } from './progress-state.component';

describe('ProgressStateComponent', () => {
  let component: ProgressStateComponent;
  let fixture: ComponentFixture<ProgressStateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ProgressStateComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ProgressStateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
