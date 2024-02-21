import { Component, OnInit, computed } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { initFlowbite } from 'flowbite';
import { NavbarComponent } from './components/navbar/navbar.component';
import { JobCardComponent } from './features/listjobs/job-card/job-card.component';
import { JobsService } from './services/jobs.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, ButtonModule, NavbarComponent, JobCardComponent],
  providers: [],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent implements OnInit {
  ngOnInit(): void {
    initFlowbite();
  }
  constructor(private jobService: JobsService) {}
  data = this.jobService.getJobs().result;
  dataMapped = computed(() =>
    this.data().data.sort((a, b) => a.id.localeCompare(b.id))
  );
}
