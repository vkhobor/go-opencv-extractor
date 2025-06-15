import { Injectable } from '@angular/core';
import { injectMutation, injectQuery, injectQueryClient } from '@ngneat/query';
import { client, baseUrl } from './http/kiota';
import { undefToErr } from './http/undefToErr';

@Injectable({
    providedIn: 'root',
})
export class TestSurfService {
    #query = injectQuery();
    #mutate = injectMutation();
    #queryClient = injectQueryClient();

    getFrame(frameNum: number) {
        return this.#query({
            queryKey: ['test-surf', 'frames', frameNum] as const,
            queryFn: () => {
                return undefToErr(
                    client.testsurf
                        .frameList({
                            framenum: frameNum,
                        })
                        .then((x) => x.data)
                );
            },
        });
    }

    async getMatchApi(req: {
        framenum: number;
        ratiocheck: number;
        minmatches: number;
        goodmatchthreshold: number;
    }) {
        return await client.testsurf.testList(req);
    }

    getFrameUrl = (frameNum: number) =>
        `${baseUrl}/testsurf/frame?framenum=${frameNum}`;

    getMaxFrame() {
        return this.#query({
            queryKey: ['test-surf', 'frames'] as const,
            queryFn: () => {
                return undefToErr(
                    client.testsurf.testsurfList().then((x) => x.data.maxframe)
                );
            },
        });
    }

    addVideo = this.#mutate({
        mutationFn: (video: File) => {
            return client.testsurf.videoCreate({
                video,
            });
        },
        onSuccess: () => {
            this.#queryClient.invalidateQueries({
                queryKey: ['test-surf', 'frames'],
            });
        },
    });

    addReference = this.#mutate({
        mutationFn: (video: File) => {
            return client.testsurf.referenceCreate({
                video,
            });
        },
        onSuccess: () => {
            this.#queryClient.invalidateQueries({
                queryKey: ['test-surf', 'match'],
            });
        },
    });
}
