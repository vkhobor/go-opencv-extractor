import { Injectable } from '@angular/core';
import { injectMutation, injectQuery, injectQueryClient } from '@ngneat/query';
import { client } from './http/kiota';
import { CreateJob, ListJobResponse } from '../../api/models';
import { undefToErr } from './http/undefToErr';

@Injectable({
  providedIn: 'root',
})
export class JobsService {
  #query = injectQuery();
  #mutate = injectMutation();
  #queryClient = injectQueryClient();

  addJob = this.#mutate({
    mutationFn: (job: CreateJob) => client.api.jobs.post(job),
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
      initialData: [] as ListJobResponse[],
      queryFn: () => undefToErr(client.api.jobs.get()),
    });
  }

  getJobDetails(id: string) {
    return this.#query({
      queryKey: ['jobs', id] as const,
      refetchInterval: 5000,
      queryFn: () => undefToErr(client.api.jobs.byId(id).get()),
    });
  }

  getJobVideos(id: string) {
    return this.#query({
      queryKey: ['jobs', id, 'videos'] as const,
      refetchInterval: 5000,
      queryFn: () => undefToErr(client.api.jobs.byId(id).videos.get()),
    });
  }
}
