import { Component, ViewChild } from '@angular/core';
import { ModalLayoutComponent } from '../../../../components/modal/modal-layout/modal-layout.component';
import { ModalContainerComponent } from '../../../../components/modal/modal-container/modal-container.component';
import { Modal } from '../../../../models/Modal';
import { CreateNewJobFormComponent } from '../form/create-new-job-form.component';
import { CommonModule } from '@angular/common';
import { JobsService } from '../../../../services/jobs.service';
import { ButtonComponent } from '../../../../components/button/button.component';
import { CreateJob } from '../../../../../api/models';

@Component({
    selector: 'app-add-modal',
    standalone: true,
    imports: [
        CommonModule,
        ModalLayoutComponent,
        ModalContainerComponent,
        CreateNewJobFormComponent,
        ButtonComponent,
    ],
    templateUrl: './add-modal.component.html',
    styleUrl: './add-modal.component.css',
})
export class AddModalComponent implements Modal {
    @ViewChild('modal') modal!: ModalContainerComponent;
    @ViewChild('form') form!: CreateNewJobFormComponent;

    constructor(private jobsService: JobsService) {}

    data: CreateJob | undefined = undefined;
    valid: boolean = false;
    saveResult = this.jobsService.addJob.result;

    openModal(): void {
        this.modal.openModal();
    }

    save(): void {
        this.form.touchAll();
        if (this.data && this.form.valid) {
            this.jobsService.addJob
                .mutateAsync(this.data)
                .then(() => this.closeModal());
        }
    }

    closeModal(): void {
        this.modal.closeModal();
    }
}
