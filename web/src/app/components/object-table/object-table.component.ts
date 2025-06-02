import { Component, Input, computed, signal } from '@angular/core';

@Component({
    selector: 'app-object-table',
    imports: [],
    templateUrl: './object-table.component.html',
    styleUrl: './object-table.component.css'
})
export class ObjectTableComponent {
    data = signal<Record<string, any>>({});
    @Input()
    set value(value: Record<string, any>) {
        this.data.set(value);
    }

    rows = computed(() => {
        return Object.entries(this.data());
    });
}
