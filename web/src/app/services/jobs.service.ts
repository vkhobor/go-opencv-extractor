import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { injectMutation, injectQuery, injectQueryClient } from '@ngneat/query';
import { Job } from '../models/Job';
import { CreateJob } from '../models/CreateJob';
import { delay, of } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class JobsService {
  #http = inject(HttpClient);
  #query = injectQuery();
  #mutate = injectMutation();
  #queryClient = injectQueryClient();

  addJob = this.#mutate({
    mutationFn: (job: CreateJob) =>
      this.#http.post('http://localhost:3010/jobs', job),
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
      initialData: [
        {
          search_query: 'test',
          id: 'test',
          limit: 0,
          progress: {
            imported: 0,
            downloaded: 0,
            scraped: 0,
          },
        },
      ] as Job[],

      queryFn: () => {
        return this.#http.get<Job[]>('http://localhost:3010/jobs');
      },
    });
  }
}
