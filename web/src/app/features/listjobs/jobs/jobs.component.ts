import { Component, computed } from '@angular/core';
import { JobsService } from '../../../services/jobs.service';
import { JobCardComponent } from '../job-card/job-card.component';
import { RouterLink } from '@angular/router';

@Component({
  selector: 'app-jobs',
  standalone: true,
  imports: [JobCardComponent, RouterLink],
  templateUrl: './jobs.component.html',
  styleUrl: './jobs.component.css',
})
export class JobsComponent {
  constructor(private jobService: JobsService) {}
  data = this.jobService.getJobs().result;
  dataMapped = computed(() => {
    const sorted = this.data().data.sort((a, b) => a.id.localeCompare(b.id));

    return sorted.map((job) => ({
      ...job,
      progress_simple:
        (job.progress.downloaded +
          job.progress.imported +
          job.progress.scraped) /
        (3 * job.limit),
    }));
  });
}
