import {
    Component,
    ViewChild,
    computed,
    signal,
    OnInit,
    effect,
} from '@angular/core';
import { Modal } from '../../../../models/Modal';
import { ModalContainerComponent } from '../../../../components/modal/modal-container/modal-container.component';

import { ModalLayoutComponent } from '../../../../components/modal/modal-layout/modal-layout.component';
import { CreateNewJobFormComponent } from '../../../newjob/components/form/create-new-job-form.component';
import { ButtonComponent } from '../../../../components/button/button.component';
import { ReferencesService } from '../../../../services/references.service';
import enviroment from '../../../../../environments/environment';
import {
    FormControl,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';

@Component({
    selector: 'app-settings-modal',
    imports: [
        ModalLayoutComponent,
        ModalContainerComponent,
        CreateNewJobFormComponent,
        ButtonComponent,
        ReactiveFormsModule,
    ],
    templateUrl: './settings-modal.component.html',
    styleUrl: './settings-modal.component.css',
})
export class SettingsModalComponent implements Modal, OnInit {
    @ViewChild('modal') modal!: ModalContainerComponent;

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
        mseSkip: new FormControl(0.2, {
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

    filesSignal = signal<File[] | null>(null);
    fileSelected($event: any) {
        this.filesSignal.set($event.target.files as File[]);
    }

    uploadResult = this.references.upload.result;

    referencesUrls = computed(() =>
        this.currentReference().data?.BlobIds?.map(
            (r) => `${enviroment.api}/files/${r}`
        )
    );

    currentReference = this.references.getReferenceById(
        'does not matter currently'
    ).result;

    constructor(private references: ReferencesService) {}

    ngOnInit(): void {
        this.resetForm();
    }

    openModal(): void {
        this.modal.openModal();
    }

    params = computed(() => {
        const reference = this.currentReference();
        return {
            minSURFMatches: reference.data?.Minsurfmatches || 3,
            minThresholdForSURFMatches:
                reference.data?.Minthresholdforsurfmatches || 0.3,
            mseSkip: reference.data?.Mseskip || 0.2,
            ratioTestThreshold: reference.data?.Ratiotestthreshold || 0.5,
        };
    });

    serverDoesNotHaveValue = computed(() => {
        const reference = this.currentReference();

        if (!reference.data?.Minsurfmatches) return true;
        return false;
    });

    myEffect = effect(() => {
        this.params();

        this.resetForm();
    });

    save(): void {
        if (this.form.valid == false)
            throw new Error('Form attempt to save in invalid state');

        this.references.upload.mutateAsync({
            files: this.filesSignal()!,
            minSURFMatches: this.form.value.minSURFMatches!,
            minThresholdForSURFMatches:
                this.form.value.minThresholdForSURFMatches!,
            mseSkip: this.form.value.mseSkip!,
            ratioTestThreshold: this.form.value.ratioTestThreshold!,
        });
    }

    resetForm(): void {
        this.form.reset();
        this.form.setValue({
            minSURFMatches: this.params().minSURFMatches,
            minThresholdForSURFMatches:
                this.params().minThresholdForSURFMatches,
            mseSkip: this.params().mseSkip,
            ratioTestThreshold: this.params().ratioTestThreshold,
        });
    }

    closeModal(): void {
        this.resetForm();
        this.modal.closeModal();
    }
}
