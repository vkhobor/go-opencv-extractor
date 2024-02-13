import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DoneStateComponent } from './done-state.component';

describe('DoneStateComponent', () => {
  let component: DoneStateComponent;
  let fixture: ComponentFixture<DoneStateComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DoneStateComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(DoneStateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
