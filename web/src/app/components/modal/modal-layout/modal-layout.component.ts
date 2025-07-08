import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
    selector: 'app-modal-layout',
    imports: [],
    templateUrl: './modal-layout.component.html',
    styleUrl: './modal-layout.component.css'
})
export class ModalLayoutComponent {
    @Input() title = 'Modal Title';
    @Input() footer = 'false';
    @Output() closeModal = new EventEmitter();
}
