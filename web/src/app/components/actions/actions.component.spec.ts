import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ActionsComponent } from './actions.component';

describe('ActionsComponent', () => {
    let component: ActionsComponent;
    let fixture: ComponentFixture<ActionsComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [ActionsComponent],
        }).compileComponents();

        fixture = TestBed.createComponent(ActionsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
