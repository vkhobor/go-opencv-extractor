import { Component, Input } from '@angular/core';
import { ModalContainerComponent } from '../modal/modal-container/modal-container.component';
import { AddModalComponent } from '../../features/newjob/components/modal/add-modal.component';
import { SettingsModalComponent } from '../../features/settings/components/settings-modal/settings-modal.component';

@Component({
  selector: 'app-settings-toggle',
  standalone: true,
  imports: [ModalContainerComponent, AddModalComponent, SettingsModalComponent],
  templateUrl: './settings-toggle.component.html',
  styleUrl: './settings-toggle.component.css',
})
export class SettingsToggleComponent {
  @Input() id = 'default';
}
