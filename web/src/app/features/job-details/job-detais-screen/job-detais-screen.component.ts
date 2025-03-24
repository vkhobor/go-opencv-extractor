import { Component, Signal, computed, inject, signal } from '@angular/core';
import { LayoutComponent } from '../../../components/layout/layout.component';
import { JobsService } from '../../../services/jobs.service';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { JsonPipe, NgClass } from '@angular/common';
import { YoutubeEmbedComponent } from '../youtube-embed/youtube-embed.component';
import { ObjectTableComponent } from '../../../components/object-table/object-table.component';
import { ProgressComponent } from '../progress/progress.component';
import { ButtonComponent } from '../../../components/button/button.component';
import { LinkIconComponent } from '../../../components/icons/link-icon/link-icon.component';
import { JobVideo } from '../../../../api/models';

@Component({
    selector: 'app-job-detais-screen',
    standalone: true,
    imports: [
        LayoutComponent,
        JsonPipe,
        LayoutComponent,
        RouterLink,
        NgClass,
        LinkIconComponent,
        YoutubeEmbedComponent,
        ObjectTableComponent,
        ProgressComponent,
        ButtonComponent,
    ],
    templateUrl: './job-detais-screen.component.html',
    styleUrl: './job-detais-screen.component.css',
})
export class JobDetaisScreenComponent {
    private jobService = inject(JobsService);
    private route = inject(ActivatedRoute);

    details = this.jobService.getJobDetails(this.route.snapshot.params['id'])
        .result;

    // progress = this.jobService.getJobProgress(this.route.snapshot.params['id'])
    //   .result;

    videosOrig = this.jobService.getJobVideos(this.route.snapshot.params['id'])
        .result;
    videos = computed(() => {
        return this.videosOrig().data?.videos?.sort((a, b) => {
            if (
                a.downloadStatus === 'success' &&
                b.downloadStatus !== 'success'
            ) {
                return -1;
            }
            if (
                a.downloadStatus !== 'success' &&
                b.downloadStatus === 'success'
            ) {
                return 1;
            }
            return 0;
        });
    });

    openProgress = signal<boolean>(false);
    toggleProgress() {
        this.openProgress.update((prev) => !prev);
    }
    toggleVideos() {
        this.openVideos.update((prev) => !prev);
    }
    openVideos = signal<boolean>(false);

    restartPipeline(): void {
        this.jobService.restartJob.mutateAsync(this.details().data?.id!);
    }

    restartResult = this.jobService.restartJob.result;

    increaseLimit() {
        this.jobService.updateJobLimit.mutateAsync({
            id: this.details().data?.id!,
            value: this.details().data?.videoTarget! + 1,
        });
    }
}
