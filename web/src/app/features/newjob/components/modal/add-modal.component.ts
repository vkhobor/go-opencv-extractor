import { Component, computed, signal, ViewChild } from '@angular/core';
import { ModalLayoutComponent } from '../../../../components/modal/modal-layout/modal-layout.component';
import { ModalContainerComponent } from '../../../../components/modal/modal-container/modal-container.component';
import { Modal } from '../../../../models/Modal';
import { CreateNewJobFormComponent } from '../form/create-new-job-form.component';
import { CommonModule } from '@angular/common';
import { JobsService } from '../../../../services/jobs.service';
import { ButtonComponent } from '../../../../components/button/button.component';
import { CreateJob } from '../../../../../api/models';
import { SimpleFormComponent } from '../simple-form/simple-form.component';

@Component({
    selector: 'app-add-modal',
    standalone: true,
    imports: [
        CommonModule,
        ModalLayoutComponent,
        ModalContainerComponent,
        CreateNewJobFormComponent,
        SimpleFormComponent,
        ButtonComponent,
    ],
    templateUrl: './add-modal.component.html',
    styleUrl: './add-modal.component.css',
})
export class AddModalComponent implements Modal {
    @ViewChild('modal') modal!: ModalContainerComponent;
    @ViewChild('form') form!: CreateNewJobFormComponent | SimpleFormComponent;

    constructor(private jobsService: JobsService) {}

    data: CreateJob | undefined = undefined;
    valid: boolean = false;
    saveResult = this.jobsService.addJob.result;
    type = signal<'query' | 'simple'>('query');

    switchButtonText = computed(() => {
        return this.type() === 'query' ? 'Switch to Simple' : 'Switch to Query';
    });

    title = computed(() => {
        return this.type() === 'query' ? 'Add Query Job' : 'Add Simple Job';
    });

    onSwitchButtonClick(): void {
        this.form.reset();
        this.type.set(this.type() === 'query' ? 'simple' : 'query');
    }

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
