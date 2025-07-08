import { CommonModule } from '@angular/common';
import { Component, Input, computed, signal } from '@angular/core';

@Component({
    selector: 'app-button',
    imports: [CommonModule],
    templateUrl: './button.component.html',
    styleUrl: './button.component.css'
})
export class ButtonComponent {
    typeSignal = signal<'primary' | 'secondary'>('primary');
    disabledSignal = signal(false);
    sizeSignal = signal<'fit' | 'full'>('fit');

    @Input()
    public set type(value: 'primary' | 'secondary') {
        this.typeSignal.set(value);
    }
    public get type(): 'primary' | 'secondary' {
        return this.typeSignal();
    }

    @Input()
    public set size(value: 'fit' | 'full') {
        this.sizeSignal.set(value);
    }
    public get size(): 'fit' | 'full' {
        return this.sizeSignal();
    }

    @Input()
    public set disabled(value: boolean) {
        this.disabledSignal.set(value);
    }
    public get disabled(): boolean {
        return this.disabledSignal();
    }

    @Input() isLoading: boolean = false;

    private readonly variantClasses = {
        primary:
            'text-white bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 me-2 mb-2 dark:bg-blue-600 dark:hover:bg-blue-700 focus:outline-none dark:focus:ring-blue-800',
        secondary:
            'py-2.5 px-5 me-2 mb-2 text-sm font-medium text-gray-900 focus:outline-none bg-white rounded-lg border border-gray-200 hover:bg-gray-100 hover:text-blue-700 focus:z-10 focus:ring-4 focus:ring-gray-100 dark:focus:ring-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:border-gray-600 dark:hover:text-white dark:hover:bg-gray-700',
    };

    currentClass = computed(
        () =>
            this.variantClasses[this.typeSignal()] +
            ' ' +
            (this.disabledSignal() ? 'cursor-not-allowed' : '') +
            ' ' +
            (this.sizeSignal() === 'fit' ? '' : 'w-full')
    );
}
