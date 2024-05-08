import { Routes } from '@angular/router';
import { JobsComponent } from './features/listjobs/jobs/jobs.component';
import { StatsScreenComponent } from './features/statistics/stats-screen/stats-screen.component';
import { FoundImagesScreenComponent } from './features/list-found-images/found-images-screen/found-images-screen.component';
import { JobDetaisScreenComponent } from './features/job-details/job-detais-screen/job-detais-screen.component';

export const routes: Routes = [
  { path: 'jobs', component: JobsComponent },
  { path: '', redirectTo: '/jobs', pathMatch: 'full' },
  { path: 'stats', component: StatsScreenComponent },
  { path: 'jobs/:id', component: JobDetaisScreenComponent },
  { path: 'images-found', component: FoundImagesScreenComponent },
];
