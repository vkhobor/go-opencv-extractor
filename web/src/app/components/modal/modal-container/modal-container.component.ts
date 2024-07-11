import { CommonModule } from '@angular/common';
import {
    Component,
    ElementRef,
    ViewChild,
    effect,
    signal,
} from '@angular/core';

@Component({
    selector: 'app-modal-container',
    standalone: true,
    imports: [CommonModule],
    templateUrl: './modal-container.component.html',
    styleUrl: './modal-container.component.scss',
})
export class ModalContainerComponent {
    @ViewChild('dialog') dialog!: ElementRef<HTMLDialogElement>;

    opened = signal(false);

    _ = effect(() => {
        if (this.opened()) {
            this.dialog.nativeElement.showModal();
            this.dialog.nativeElement.classList.add('opened');
        } else {
            this.dialog.nativeElement.close();
            this.dialog.nativeElement.classList.remove('opened');
        }
    });

    openModal() {
        this.opened.set(true);
    }

    closeModal() {
        this.opened.set(false);
    }

    ngAfterViewInit() {
        this.dialog.nativeElement.addEventListener(
            'click',
            (event: MouseEvent) => {
                const target = event.target as Element;
                if (target.nodeName === 'DIALOG') {
                    this.opened.set(false);
                }
            }
        );
    }
}
