import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { injectQuery } from '@ngneat/query';
import { Job } from '../models/Job';

@Injectable({
  providedIn: 'root',
})
export class JobsService {
  #http = inject(HttpClient);
  #query = injectQuery();
  getJobs() {
    return this.#query({
      queryKey: ['jobss'] as const,
      refetchInterval: 5000,
      initialData: [],

      queryFn: () => {
        return this.#http.get<
          Job[]
        >('http://localhost:3010/jobs');
      },
    });
  }
}
