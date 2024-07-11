import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FoundImagesScreenComponent } from './found-images-screen.component';

describe('FoundImagesScreenComponent', () => {
    let component: FoundImagesScreenComponent;
    let fixture: ComponentFixture<FoundImagesScreenComponent>;

    beforeEach(async () => {
        await TestBed.configureTestingModule({
            imports: [FoundImagesScreenComponent],
        }).compileComponents();

        fixture = TestBed.createComponent(FoundImagesScreenComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
    });

    it('should create', () => {
        expect(component).toBeTruthy();
    });
});
