import { Component, Signal, computed, inject, signal } from '@angular/core';
import { LayoutComponent } from '../../../components/layout/layout.component';
import { JobsService } from '../../../services/jobs.service';
import { ActivatedRoute } from '@angular/router';
import { toObservable, toSignal } from '@angular/core/rxjs-interop';
import { filter, startWith, switchMap } from 'rxjs';
import {
  DefinedQueryObserverResult,
  QueryObserverLoadingErrorResult,
  QueryObserverResult,
  createPendingObserverResult,
} from '@ngneat/query';
import { JobDetails } from '../../../models/JobDetails';
import { JsonPipe, NgClass } from '@angular/common';
import { YoutubeEmbedComponent } from '../youtube-embed/youtube-embed.component';
import { ObjectTableComponent } from '../../../components/object-table/object-table.component';
import { ProgressComponent } from '../progress/progress.component';

@Component({
  selector: 'app-job-detais-screen',
  standalone: true,
  imports: [
    LayoutComponent,
    JsonPipe,
    NgClass,
    YoutubeEmbedComponent,
    ObjectTableComponent,
    ProgressComponent,
  ],
  templateUrl: './job-detais-screen.component.html',
  styleUrl: './job-detais-screen.component.css',
})
export class JobDetaisScreenComponent {
  private jobService = inject(JobsService);
  private route = inject(ActivatedRoute);

  details = this.jobService.getJobDetails(this.route.snapshot.params['id'])
    .result;

  // progress = this.jobService.getJobProgress(this.route.snapshot.params['id'])
  //   .result;

  videos = this.jobService.getJobVideos(this.route.snapshot.params['id'])
    .result;

  openProgress = signal<boolean>(false);
  toggleProgress() {
    this.openProgress.update((prev) => !prev);
  }
  toggleVideos() {
    this.openVideos.update((prev) => !prev);
  }
  openVideos = signal<boolean>(false);
}
