import { Component, OnInit, computed } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { CardModule } from 'primeng/card';
import { AccordionModule } from 'primeng/accordion';
import { ProgressBarModule } from 'primeng/progressbar';
import { ToastModule } from 'primeng/toast';
import { ToolbarModule } from 'primeng/toolbar';
import { initFlowbite } from 'flowbite';
import { NavbarComponent } from './components/navbar/navbar.component';
import { JobCardComponent } from './components/job-card/job-card.component';
import { AddSearchModalComponent } from './components/add-search-modal/add-search-modal.component';
import { AddNewService } from './services/add-new.service';
import { JobsService } from './services/jobs.service';
import { HttpClient } from '@angular/common/http';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet,
    ButtonModule,
    NavbarComponent,
    JobCardComponent,
    AddSearchModalComponent,
  ],
  providers: [AddNewService],
  templateUrl: './app.component.html',
  styleUrl: './app.component.css',
})
export class AppComponent implements OnInit {
  ngOnInit(): void {
    initFlowbite();
  }
  constructor(
    public addNewService: AddNewService,
    private jobService: JobsService
  ) {}
  data = this.jobService.getJobs().result;
  dataMapped = computed(() =>
    this.data()
      .data.sort((a, b) => a.id.localeCompare(b.id))
      .map((job) => {
        return {
          id: job.id,
          search: job.search_query,
          progress: job.progress.completed / job.progress.total,
          done: job.progress.completed === job.progress.total,
          amount: job.limit,
        };
      })
  );
}
