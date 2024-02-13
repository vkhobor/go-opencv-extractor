import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-progress-state',
  standalone: true,
  imports: [],
  templateUrl: './progress-state.component.html',
  styleUrl: './progress-state.component.css',
})
export class ProgressStateComponent {
  @Input() job!: { search: string; progress: number; amount: number };
}
