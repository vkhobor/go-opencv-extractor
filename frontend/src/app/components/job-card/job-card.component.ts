import { Component, Input } from '@angular/core';
import { ProgressStateComponent } from './components/progress-state/progress-state.component';
import { DoneStateComponent } from './components/done-state/done-state.component';

@Component({
  selector: 'app-job-card',
  standalone: true,
  imports: [ProgressStateComponent, DoneStateComponent],
  templateUrl: './job-card.component.html',
  styleUrl: './job-card.component.css',
})
export class JobCardComponent {
  @Input() job!: {
    search: string;
    progress: number;
    amount: number;
    done: boolean;
  };
}
