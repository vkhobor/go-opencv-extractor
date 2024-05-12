import { Component, EventEmitter, Input, Output, signal } from '@angular/core';
import { Flowbite } from '../../util/flowbiteFix';

export interface Action {
  id: string;
  label: string;
}

@Component({
  selector: 'app-actions',
  standalone: true,
  imports: [],
  templateUrl: './actions.component.html',
  styleUrl: './actions.component.css',
})
@Flowbite()
export class ActionsComponent {
  @Input() id = this.generateUniqueString();
  @Input() actions: Action[] = [];
  @Output() actionSelected = new EventEmitter<Action>();

  open = signal(false);
  toggleOpen(): void {
    this.open.update((prev) => !prev);
  }

  generateUniqueString() {
    return Math.random().toString(36).substring(2) + Date.now().toString(36);
  }
}
