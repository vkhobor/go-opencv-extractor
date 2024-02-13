import { Component, Input } from '@angular/core';
import { AddSearchModalComponent } from '../add-search-modal/add-search-modal.component';
import { AddNewService } from '../../services/add-new.service';

@Component({
  selector: 'app-settings-toggle',
  standalone: true,
  imports: [AddSearchModalComponent],
  templateUrl: './settings-toggle.component.html',
  styleUrl: './settings-toggle.component.css',
})
export class SettingsToggleComponent {
  @Input() id = 'default';
  constructor(public addNewService: AddNewService) {}
}
