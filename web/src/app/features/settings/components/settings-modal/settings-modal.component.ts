import { Component, ViewChild, computed, signal } from '@angular/core';
import { Modal } from '../../../../models/Modal';
import { ModalContainerComponent } from '../../../../components/modal/modal-container/modal-container.component';
import { CommonModule } from '@angular/common';
import { ModalLayoutComponent } from '../../../../components/modal/modal-layout/modal-layout.component';
import { CreateNewJobFormComponent } from '../../../newjob/components/form/create-new-job-form.component';
import { ButtonComponent } from '../../../../components/button/button.component';
import { ReferencesService } from '../../../../services/references.service';

@Component({
  selector: 'app-settings-modal',
  standalone: true,
  imports: [
    CommonModule,
    ModalLayoutComponent,
    ModalContainerComponent,
    CreateNewJobFormComponent,
    ButtonComponent,
  ],
  templateUrl: './settings-modal.component.html',
  styleUrl: './settings-modal.component.css',
})
export class SettingsModalComponent implements Modal {
  @ViewChild('modal') modal!: ModalContainerComponent;

  filesSignal = signal<File[] | null>(null);
  fileSelected($event: any) {
    this.filesSignal.set($event.target.files as File[]);
  }

  uploadResult = this.references.upload.result;
  deleteAllResult = this.references.deleteAll.result;

  referencesValues = this.references.getReferences().result;
  referencesUrls = computed(() =>
    this.referencesValues().data?.map((r) => `/files/${r.id}`)
  );

  constructor(private references: ReferencesService) {}

  openModal(): void {
    this.modal.openModal();
  }

  save(): void {
    this.references.upload.mutateAsync(this.filesSignal()!);
  }

  deleteAll(): void {
    this.references.deleteAll.mutateAsync({});
  }

  closeModal(): void {
    this.modal.closeModal();
  }
}
