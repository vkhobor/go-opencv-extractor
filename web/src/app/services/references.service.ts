import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { injectMutation, injectQuery, injectQueryClient } from '@ngneat/query';
import { DefaultHttpProxyService } from './http/default-http-proxy.service';
import { Observable } from 'rxjs';
import { Job } from '../models/Job';

@Injectable({
  providedIn: 'root',
})
export class ReferencesService {
  #http = inject(DefaultHttpProxyService);
  #query = injectQuery();
  #mutate = injectMutation();
  #queryClient = injectQueryClient();

  upload = this.#mutate({
    mutationFn: (files: File[]) => {
      const formData = new FormData();
      for (const file of files) {
        formData.append(file.name, file);
      }
      return this.#http.post('/references', formData);
    },

    onSuccess: () => {
      this.#queryClient.invalidateQueries({
        queryKey: ['references'],
      });
    },
  });

  getReferences() {
    return this.#query({
      queryKey: ['references'] as const,
      queryFn: () => {
        return this.#http.get('/references') as Observable<{ id: string }[]>;
      },
    });
  }

  constructor() {}
}
