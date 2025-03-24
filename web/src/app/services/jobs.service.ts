import { Injectable } from '@angular/core';
import { injectMutation, injectQuery, injectQueryClient } from '@ngneat/query';
import { client } from './http/kiota';
import { CreateJob, ListJobBody } from '../../api/models';
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

    restartJob = this.#mutate({
        mutationFn: (id: string) =>
            client.api.jobs.byId(id).actions.restart.post(),
        onSuccess: () => {
            this.#queryClient.invalidateQueries({
                queryKey: ['jobs'],
            });
        },
    });

    updateJobLimit = this.#mutate({
        mutationFn: ({ id, value }: { id: string; value: number }) =>
            client.api.jobs.byId(id).actions.updateLimit.post({ limit: value }),
        onSuccess: (_, { id }) => {
            this.#queryClient.invalidateQueries({
                queryKey: ['jobs', id],
            });
        },
    });

    getJobs() {
        return this.#query({
            queryKey: ['jobs'] as const,
            refetchInterval: 5000,
            initialData: [] as ListJobBody[],
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
