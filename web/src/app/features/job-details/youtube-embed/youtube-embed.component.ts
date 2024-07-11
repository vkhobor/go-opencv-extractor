import { Component, Input } from '@angular/core';
import { SafePipe } from '../../../pipes/safe.pipe';

@Component({
    selector: 'app-youtube-embed',
    standalone: true,
    imports: [SafePipe],
    templateUrl: './youtube-embed.component.html',
    styleUrl: './youtube-embed.component.css',
})
export class YoutubeEmbedComponent {
    @Input() id: string = undefined!;

    get embedUrl() {
        return `http://img.youtube.com/vi/${this.id}/hqdefault.jpg`;
    }

    get link() {
        return `https://www.youtube.com/watch?v=${this.id}`;
    }
}
