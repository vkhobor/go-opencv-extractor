import { Injectable } from '@angular/core';
import { injectMutation, injectQuery, injectQueryClient } from '@ngneat/query';
import { client } from './http/kiota';
import { CreateJob, ListJobBody, ListVideoBody } from '../../api/Api';
import { undefToErr } from './http/undefToErr';

@Injectable({
    providedIn: 'root',
})
export class JobsService {
    #query = injectQuery();
    #mutate = injectMutation();
    #queryClient = injectQueryClient();

    // TODO remove
    addJob = this.#mutate({
        mutationFn: (job: CreateJob) => {
            return new Promise(() => {});
        },
        onSuccess: () => {
            this.#queryClient.invalidateQueries({
                queryKey: ['jobs'],
            });
        },
    });

    addVideoJob = this.#mutate({
        mutationFn: (job: { blob: File; filterId: string; name: string }) => {
            return client.api.jobsVideoCreate({
                file: job.blob,
                filter_id: job.filterId,
                name: job.name,
            });
        },
        onSuccess: () => {
            this.#queryClient.invalidateQueries({
                queryKey: ['jobs'],
            });
        },
    });

    restartJob = this.#mutate({
        mutationFn: (id: string) => new Promise(() => {}),
        onSuccess: () => {
            this.#queryClient.invalidateQueries({
                queryKey: ['jobs'],
            });
        },
    });

    updateJobLimit = this.#mutate({
        mutationFn: ({ id, value }: { id: string; value: number }) =>
            new Promise(() => {}),
        onSuccess: (_, { id }) => {
            this.#queryClient.invalidateQueries({
                queryKey: ['jobs', id],
            });
        },
    });

    getVideos() {
        return this.#query({
            queryKey: ['jobs'] as const,
            refetchInterval: 5000,
            initialData: [] as ListVideoBody[],
            queryFn: () =>
                undefToErr(client.api.videosList().then((x) => x.data)),
        });
    }

    getJobDetails(id: string) {
        return this.#query({
            queryKey: ['jobs', id] as const,
            refetchInterval: 5000,
            queryFn: () => undefToErr(new Promise(() => {})),
        });
    }

    getJobVideos(id: string) {
        return this.#query({
            queryKey: ['jobs', id, 'videos'] as const,
            refetchInterval: 5000,
            queryFn: () => undefToErr(new Promise(() => {})),
        });
    }
}
