import { Component, Input } from '@angular/core';
import { Job } from '../../../models/Job';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-job-card',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './job-card.component.html',
  styleUrl: './job-card.component.css',
})
export class JobCardComponent {
  @Input() job!: Job;
  open = false;
}
