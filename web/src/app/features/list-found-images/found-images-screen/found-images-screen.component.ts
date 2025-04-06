import { Component, computed, effect, inject, signal } from '@angular/core';
import { LayoutComponent } from '../../../components/layout/layout.component';
import { ImagesService } from '../../../services/images.service';
import { JsonPipe } from '@angular/common';
import enviroment from '../../../../enviroments/enviroment';
import { ZipService } from '../../../services/zip.service';
import { ActionsComponent } from '../../../components/actions/actions.component';
import { ActivatedRoute, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { BadgeComponent } from '../../../components/badge/badge.component';

@Component({
    selector: 'app-found-images-screen',
    standalone: true,
    imports: [LayoutComponent, BadgeComponent, JsonPipe, ActionsComponent],
    templateUrl: './found-images-screen.component.html',
    styleUrl: './found-images-screen.component.css',
})
export class FoundImagesScreenComponent {
    imagesService = inject(ImagesService);
    zipService = inject(ZipService);
    activatedRoute = inject(ActivatedRoute);
    router = inject(Router);

    activatedRouteParams = toSignal(this.activatedRoute.queryParams);
    __ = effect(() => {
        console.log(this.imagePage());
    });
    requestParams = computed<{ youtubeId: string }>(() => ({
        youtubeId:
            this.activatedRouteParams() !== undefined
                ? this.activatedRouteParams()!['youtube_id']
                : undefined,
    }));

    readonly pageSize = 10;
    currentPageNumber = signal(0);

    actions = [{ id: 'exportAll', label: 'Export all' }];

    onActionSelected(action: { id: string }) {
        switch (action.id) {
            case 'exportAll':
                this.exportAll();
                break;
        }
    }

    imagePageQuery = this.imagesService.getImagePage(
        this.currentPageNumber(),
        this.pageSize,
        this.requestParams().youtubeId
    );
    imagePage = this.imagePageQuery.result;

    _ = effect(() => {
        console.log('updateQuery', this.requestParams());
        this.imagePageQuery.updateOptions({
            queryKey: [
                'images',
                this.requestParams(),
                this.currentPageNumber(),
                this.pageSize,
            ] as const,
            enabled: true,
            queryFn: () =>
                this.imagesService.getImagePageApi(
                    this.currentPageNumber(),
                    this.pageSize,
                    this.requestParams().youtubeId
                ),
        });
    });

    referencesUrls = computed(
        () =>
            this.imagePage().data?.pictures?.map(
                (r) => `${enviroment.api}/files/${r.blob_id}`
            ) ?? []
    );

    filters = computed(() => [
        {
            label: 'Youtube ID',
            value: this.requestParams().youtubeId,
            enabled: this.requestParams().youtubeId !== undefined,
            onDismiss: () => {
                this.router.navigate([], {
                    queryParams: {
                        youtube_id: null,
                    },
                    queryParamsHandling: 'merge',
                });
            },
        },
    ]);

    filtersEnabled = computed(() => this.filters().filter((f) => f.enabled));

    exportAll() {
        this.zipService.downloadZip();
    }

    previous() {
        if (this.currentPageNumber() > 0) {
            this.currentPageNumber.update((prev) => prev - 1);
        }
    }

    next() {
        if (
            this.imagePage().data!.total! >
            (this.currentPageNumber() + 1) * this.pageSize
        ) {
            this.currentPageNumber.update((prev) => prev + 1);
        }
    }
}
