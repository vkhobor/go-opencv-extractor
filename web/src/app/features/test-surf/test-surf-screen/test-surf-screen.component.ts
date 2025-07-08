import { Component, computed, effect, inject, signal } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { LayoutComponent } from '../../../components/layout/layout.component';
import {
    FormControl,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { resource } from '@angular/core';
import { TestSurfService } from '../../../services/test-surf.service';
import { startWith } from 'rxjs';

@Component({
    selector: 'app-test-surf-screen',
    imports: [LayoutComponent, ReactiveFormsModule],
    templateUrl: './test-surf-screen.component.html',
    styleUrl: './test-surf-screen.component.css',
})
export class TestSurfScreenComponent {
    testSurfService = inject(TestSurfService);

    selectedFrame = signal(1);
    videoUploaded = signal(false);
    referenceUploaded = signal(false);

    maxFrameResult = this.testSurfService.getMaxFrame().result;

    maxFrame = computed(() => {
        const data = this.maxFrameResult();
        if (data.isSuccess) {
            return data.data;
        }
        return 100;
    });

    frameUrl = computed(() =>
        this.testSurfService.getFrameUrl(this.selectedFrame())
    );

    frameSelected(event: Event) {
        const tEvent = event.target as HTMLInputElement;
        this.selectedFrame.set(tEvent.valueAsNumber);
    }

    async testVideoSelected(event: Event) {
        const input = event.target as HTMLInputElement;

        if (input.files && input.files.length > 0) {
            const file = input.files[0];
            await this.testSurfService.addVideo.mutateAsync(file);
            this.videoUploaded.set(true);
        }
    }

    async referenceSelected(event: Event) {
        const input = event.target as HTMLInputElement;

        if (input.files && input.files.length > 0) {
            const file = input.files[0];
            await this.testSurfService.addReference.mutateAsync(file);
            this.referenceUploaded.set(true);
        }
    }

    form = new FormGroup({
        minSURFMatches: new FormControl(3, {
            nonNullable: true,
            validators: [
                Validators.required,
                Validators.min(0),
                Validators.pattern(/^\d+$/),
            ],
        }),
        minThresholdForSURFMatches: new FormControl(0.3, {
            nonNullable: true,
            validators: [
                Validators.required,
                Validators.min(0),
                Validators.max(1),
            ],
        }),

        ratioTestThreshold: new FormControl(0.5, {
            nonNullable: true,
            validators: [
                Validators.required,
                Validators.min(0),
                Validators.max(1),
            ],
        }),
    });

    formValueSignal = toSignal(
        this.form.valueChanges.pipe(startWith(this.form.value))
    );

    isMatch = resource({
        params: () => ({
            framenum: this.selectedFrame(),
            ratiocheck: this.formValueSignal()?.ratioTestThreshold,
            minmatches: this.formValueSignal()?.minSURFMatches,
            goodmatchthreshold:
                this.formValueSignal()?.minThresholdForSURFMatches,
        }),
        loader: async ({ params }) => {
            if (params.framenum === undefined) {
                return Promise.reject(new Error('Invalid input'));
            }
            if (params.ratiocheck === undefined) {
                return Promise.reject(new Error('Invalid input'));
            }
            if (params.minmatches === undefined) {
                return Promise.reject(new Error('Invalid input'));
            }
            if (params.goodmatchthreshold === undefined) {
                return Promise.reject(new Error('Invalid input'));
            }

            const valid: {
                framenum: number;
                ratiocheck: number;
                minmatches: number;
                goodmatchthreshold: number;
            } = {
                framenum: params.framenum,
                ratiocheck: params.ratiocheck,
                minmatches: params.minmatches,
                goodmatchthreshold: params.goodmatchthreshold,
            };
            try {
                const x = await this.testSurfService.getMatchApi(valid);
                return x.data.matched;
            } catch (error) {
                throw new Error('Unexpected api error' + error);
            }
        },
    });

    showErrorAlert = computed(() => {
        const status = this.isMatch.status();
        const error = this.isMatch.error();
        const videoUploaded = this.videoUploaded();
        const referenceUploaded = this.referenceUploaded();

        if (videoUploaded === false) {
            return 'Upload a video to begin with! After that you can scrub through the video.';
        }

        if (referenceUploaded === false) {
            return 'Upload a reference image to use for the matching algo.';
        }

        if (status === 'error') {
            return error?.message;
        }

        return undefined;
    });

    showLoading = computed(() => {
        const error = this.showErrorAlert();
        const status = this.isMatch.status();

        if (!error && status === 'loading') {
            return true;
        }

        return false;
    });

    showSuccessAlert = computed(() => {
        const loading = this.showLoading();
        const error = this.showErrorAlert();
        const hasValue = this.isMatch.hasValue();

        if (!loading && !error && hasValue) {
            return this.isMatch.value();
        }

        return false;
    });

    showNonMatched = computed(() => {
        const loading = this.showLoading();
        const error = this.showErrorAlert();
        const hasValue = this.isMatch.hasValue();

        if (!loading && !error && hasValue) {
            return !this.isMatch.value();
        }

        return false;
    });
}
