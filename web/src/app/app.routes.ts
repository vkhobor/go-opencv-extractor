import { Routes } from '@angular/router';
import { JobsComponent } from './features/listjobs/jobs/jobs.component';
import { StatsScreenComponent } from './features/statistics/stats-screen/stats-screen.component';

export const routes: Routes = [
  { path: '', component: JobsComponent },
  { path: 'stats', component: StatsScreenComponent },
];
