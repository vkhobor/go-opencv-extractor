import { Component, computed, EventEmitter, Output } from '@angular/core';
import { CreateJob } from '../../../../../api/models';
import {
    FormBuilder,
    FormControl,
    FormsModule,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { FilterService } from '../../../../services/filter.service';

@Component({
    selector: 'app-simple-form',
    standalone: true,
    imports: [FormsModule, ReactiveFormsModule],
    templateUrl: './simple-form.component.html',
    styleUrl: './simple-form.component.css',
})
export class SimpleFormComponent {
    @Output() valid = new EventEmitter<boolean>();
    @Output() data = new EventEmitter<CreateJob | undefined>();

    form = this.fb.group(
        {
            youtubeId: new FormControl<string | null>(null, [
                Validators.required,
                Validators.minLength(5),
            ]),
            filter: new FormControl<string | null>('', [
                Validators.required,
                Validators.minLength(1),
            ]),
        },
        { updateOn: 'change' }
    );

    filters = this.filterService.getFilters().result;
    filterOptions = computed(() =>
        this.filters().data?.map((f) => ({ label: f.name, value: f.id }))
    );

    constructor(
        private fb: FormBuilder,
        private filterService: FilterService
    ) {}

    public get youtubeId() {
        return this.form.get('youtubeId');
    }

    touchAll() {
        this.form.markAllAsTouched();
    }

    reset() {
        this.form.reset();
    }

    ngOnInit() {
        this.form.valueChanges.subscribe((data) => {
            this.valid.emit(this.form.valid);

            if (this.form.valid) {
                this.data.emit({
                    isSingleVideo: true,
                    isQueryBased: false,
                    youtubeId: data.youtubeId!,
                    filterId: data.filter!,
                });
            } else {
                this.data.emit(undefined);
            }
        });
    }
}
