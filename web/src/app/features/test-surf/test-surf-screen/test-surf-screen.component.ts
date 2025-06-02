import { Component, computed, effect, inject, signal } from '@angular/core';
import { LayoutComponent } from '../../../components/layout/layout.component';
import { CommonModule } from '@angular/common';
import { ModalLayoutComponent } from '../../../components/modal/modal-layout/modal-layout.component';
import {
    FormControl,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { resource, Signal } from '@angular/core';
import { ButtonComponent } from '../../../components/button/button.component';
import { ModalContainerComponent } from '../../../components/modal/modal-container/modal-container.component';
import { CreateNewJobFormComponent } from '../../newjob/components/form/create-new-job-form.component';
import { TestSurfService } from '../../../services/test-surf.service';

@Component({
    selector: 'app-test-surf-screen',
    standalone: true,
    imports: [LayoutComponent, CommonModule, ReactiveFormsModule],
    templateUrl: './test-surf-screen.component.html',
    styleUrl: './test-surf-screen.component.css',
})
export class TestSurfScreenComponent {
    testSurfService = inject(TestSurfService);

    selectedFrame = signal(1);
    maxFrame = this.testSurfService.getMaxFrame().result;

    frameUrl = computed(() =>
        this.testSurfService.getFrameUrl(this.selectedFrame())
    );

    // isMatch = computedAsync(() => this.usersService.getUser(+this.id()).result$, {
    //     initialValue: createPendingObserverResult<User>(),
    //   });

    _ = effect(() => console.log(this.frameUrl()));

    frameSelected(event: Event) {
        const tEvent = event.target as HTMLInputElement;
        this.selectedFrame.set(tEvent.valueAsNumber);
    }

    testVideoSelected(event: Event) {
        const input = event.target as HTMLInputElement;

        if (input.files && input.files.length > 0) {
            const file = input.files[0];
            this.testSurfService.addVideo.mutate(file);
        }
    }

    referenceSelected(event: Event) {
        const input = event.target as HTMLInputElement;

        if (input.files && input.files.length > 0) {
            const file = input.files[0];
            this.testSurfService.addReference.mutate(file);
        }
    }

    form = new FormGroup({
        minSURFMatches: new FormControl(0, {
            nonNullable: true,
            validators: [
                Validators.required,
                Validators.min(0),
                Validators.pattern(/^\d+$/),
            ],
        }),
        minThresholdForSURFMatches: new FormControl(0, {
            nonNullable: true,
            validators: [
                Validators.required,
                Validators.min(0),
                Validators.max(1),
            ],
        }),
        mseSkip: new FormControl(0, {
            nonNullable: true,
            validators: [
                Validators.required,
                Validators.min(0),
                Validators.max(1),
            ],
        }),
        ratioTestThreshold: new FormControl(0, {
            nonNullable: true,
            validators: [
                Validators.required,
                Validators.min(0),
                Validators.max(1),
            ],
        }),
    });
}
