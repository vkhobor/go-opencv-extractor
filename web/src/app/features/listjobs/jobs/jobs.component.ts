import { Component, ViewChild, computed } from '@angular/core';
import { JobsService } from '../../../services/jobs.service';
import { RouterLink } from '@angular/router';
import { LayoutComponent } from '../../../components/layout/layout.component';
import { AddModalComponent } from '../../newjob/components/modal/add-modal.component';
import { SettingsModalComponent } from '../../settings/components/settings-modal/settings-modal.component';
import {
    Action,
    ActionsComponent,
} from '../../../components/actions/actions.component';
import { initFlowbite } from 'flowbite';

@Component({
    selector: 'app-jobs',
    standalone: true,
    imports: [
        ActionsComponent,
        SettingsModalComponent,
        RouterLink,
        LayoutComponent,
        AddModalComponent,
    ],
    templateUrl: './jobs.component.html',
    styleUrl: './jobs.component.css',
})
export class JobsComponent {
    @ViewChild(AddModalComponent) addModal!: AddModalComponent;
    @ViewChild(SettingsModalComponent) settingsModal!: SettingsModalComponent;

    constructor(private jobService: JobsService) {}
    data = this.jobService.getJobs().result;

    actions: Action[] = [
        {
            id: 'Add',
            label: 'Add',
        },
        {
            id: 'Default filter',
            label: 'Default filter',
        },
    ];

    onActionSelected(action: Action) {
        switch (action.id) {
            case 'Add':
                this.addModal.openModal();
                break;
            case 'Default filter':
                this.settingsModal.openModal();
                break;
        }
    }

    dataMapped = computed(() => {
        const sorted = this.data().data.sort((a, b) =>
            a.id!.localeCompare(b.id!)
        );

        return sorted.map((job) => ({
            ...job,
        }));
    });
}
