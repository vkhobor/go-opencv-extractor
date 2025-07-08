import { Routes } from '@angular/router';
import { JobsComponent } from './features/listjobs/jobs/jobs.component';
import { FoundImagesScreenComponent } from './features/list-found-images/found-images-screen/found-images-screen.component';
import { TestSurfScreenComponent } from './features/test-surf/test-surf-screen/test-surf-screen.component';

export const routes: Routes = [
    { path: 'jobs', component: JobsComponent },
    { path: '', redirectTo: '/jobs', pathMatch: 'full' },
    { path: 'images-found', component: FoundImagesScreenComponent },
    { path: 'test-surf', component: TestSurfScreenComponent },
];
