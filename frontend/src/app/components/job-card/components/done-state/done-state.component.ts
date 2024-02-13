import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-done-state',
  standalone: true,
  imports: [],
  templateUrl: './done-state.component.html',
  styleUrl: './done-state.component.css',
})
export class DoneStateComponent {
  @Input() job!: { search: string; progress: number; amount: number };
}
