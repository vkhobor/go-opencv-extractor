import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SettingsToggleComponent } from './settings-toggle.component';

describe('SettingsToggleComponent', () => {
  let component: SettingsToggleComponent;
  let fixture: ComponentFixture<SettingsToggleComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [SettingsToggleComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(SettingsToggleComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
