import { CommonModule } from '@angular/common';
import {
    Component,
    computed,
    EventEmitter,
    Input,
    Output,
    signal,
} from '@angular/core';

@Component({
    selector: 'app-badge',
    imports: [CommonModule],
    templateUrl: './badge.component.html',
    styleUrl: './badge.component.css'
})
export class BadgeComponent {
    @Output() onDismiss = new EventEmitter<void>();

    dismiss() {
        this.onDismiss.emit();
    }
    typeSignal = signal<'primary' | 'secondary'>('primary');
    @Input()
    text = '';

    @Input()
    public set type(value: 'primary' | 'secondary') {
        this.typeSignal.set(value);
    }
    public get type(): 'primary' | 'secondary' {
        return this.typeSignal();
    }

    private readonly variantClasses = {
        primary: {
            badge: 'inline-flex items-center px-2 py-1 me-2 text-sm font-medium text-blue-800 bg-blue-100 rounded dark:bg-blue-900 dark:text-blue-300',
            dismissBtn:
                'inline-flex items-center p-1 ms-2 text-sm text-blue-400 bg-transparent rounded-sm hover:bg-blue-200 hover:text-blue-900 dark:hover:bg-blue-800 dark:hover:text-blue-300',
        },
        secondary: {
            badge: 'inline-flex items-center px-2 py-1 me-2 text-sm font-medium text-gray-800 bg-gray-100 rounded dark:bg-gray-700 dark:text-gray-300',
            dismissBtn:
                'inline-flex items-center p-1 ms-2 text-sm text-gray-400 bg-transparent rounded-sm hover:bg-gray-200 hover:text-gray-900 dark:hover:bg-gray-600 dark:hover:text-gray-300',
        },
    };

    currentClassBadge = computed(
        () => this.variantClasses[this.typeSignal()].badge
    );
    currentClassDismiss = computed(
        () => this.variantClasses[this.typeSignal()].dismissBtn
    );
}
