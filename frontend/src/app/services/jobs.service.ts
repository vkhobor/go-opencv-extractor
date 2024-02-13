import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { injectQuery } from '@ngneat/query';

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
          {
            search_query: string;
            id: string;
            limit: number;
            progress: { total: number; completed: number };
          }[]
        >('http://localhost:3010/jobs');
      },
    });
  }
}
