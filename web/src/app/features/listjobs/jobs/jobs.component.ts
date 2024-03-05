import { Component, computed } from '@angular/core';
import { JobsService } from '../../../services/jobs.service';
import { JobCardComponent } from '../job-card/job-card.component';

@Component({
  selector: 'app-jobs',
  standalone: true,
  imports: [JobCardComponent],
  templateUrl: './jobs.component.html',
  styleUrl: './jobs.component.css',
})
export class JobsComponent {
  constructor(private jobService: JobsService) {}
  data = this.jobService.getJobs().result;
  dataMapped = computed(() =>
    this.data().data.sort((a, b) => a.id.localeCompare(b.id))
  );
}
