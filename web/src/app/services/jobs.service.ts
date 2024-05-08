import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { injectMutation, injectQuery, injectQueryClient } from '@ngneat/query';
import { Job } from '../models/Job';
import { CreateJob } from '../models/CreateJob';
import { DefaultHttpProxyService } from './http/default-http-proxy.service';
import { Observable } from 'rxjs';
import { JobDetails } from '../models/JobDetails';
import { JobWithVideos } from '../models/JobWithVideos';
import { Progress } from '../models/JobProgress';

@Injectable({
  providedIn: 'root',
})
export class JobsService {
  #http = inject(DefaultHttpProxyService);
  #query = injectQuery();
  #mutate = injectMutation();
  #queryClient = injectQueryClient();

  addJob = this.#mutate({
    mutationFn: (job: CreateJob) => this.#http.post('/jobs', job),
    onSuccess: () => {
      this.#queryClient.invalidateQueries({
        queryKey: ['jobs'],
      });
    },
  });

  getJobs() {
    return this.#query({
      queryKey: ['jobs'] as const,
      refetchInterval: 5000,
      initialData: [] as Job[],

      queryFn: () => {
        return this.#http.get('/jobs') as Observable<Job[]>;
      },
    });
  }

  getJobDetails(id: string) {
    return this.#query({
      queryKey: ['jobs', id] as const,
      queryFn: () => {
        return this.#http.get(`/jobs/${id}`) as Observable<JobDetails>;
      },
    });
  }

  getJobVideos(id: string) {
    return this.#query({
      queryKey: ['jobs', id, 'videos'] as const,
      queryFn: () => {
        return this.#http.get(
          `/jobs/${id}/videos`
        ) as Observable<JobWithVideos>;
      },
    });
  }

  getJobProgress(id: string) {
    return this.#query({
      queryKey: ['jobs', id, 'progress'] as const,
      queryFn: () => {
        return this.#http.get(`/jobs/${id}/progress`) as Observable<Progress>;
      },
    });
  }
}
