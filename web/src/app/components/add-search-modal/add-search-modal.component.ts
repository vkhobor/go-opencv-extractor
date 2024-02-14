import { Component, EventEmitter, Input, Output } from '@angular/core';
import { InfoCircleIcon } from 'primeng/icons/infocircle';

@Component({
  selector: 'app-add-search-modal',
  standalone: true,
  imports: [],
  templateUrl: './add-search-modal.component.html',
  styleUrl: './add-search-modal.component.css',
})
export class AddSearchModalComponent {
  @Input() id = 'default-add-search-modal';
  @Input() show = false;
  @Output() onHide = new EventEmitter();
}
