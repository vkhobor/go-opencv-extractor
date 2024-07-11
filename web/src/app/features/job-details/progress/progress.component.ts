import { Component, Input } from '@angular/core';
import { Progress } from '../../../models/JobProgress';

@Component({
    selector: 'app-progress',
    standalone: true,
    imports: [],
    templateUrl: './progress.component.html',
    styleUrl: './progress.component.css',
})
export class ProgressComponent {
    @Input() job: Progress = undefined!;
}
