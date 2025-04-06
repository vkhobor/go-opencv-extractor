import { Component, EventEmitter, Output, computed } from '@angular/core';
import {
    FormBuilder,
    FormControl,
    FormsModule,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { FilterService } from '../../../../services/filter.service';
import { JobsService } from '../../../../services/jobs.service';

@Component({
    selector: 'app-create-new-job-form',
    standalone: true,
    imports: [FormsModule, ReactiveFormsModule],
    templateUrl: './create-new-job-form.component.html',
})
export class CreateNewJobFormComponent {
    @Output() valid = new EventEmitter<boolean>();
    @Output() data = new EventEmitter<{
        name: string;
        filterId: string;
        file: File;
    }>();

    form = this.fb.group({
        name: new FormControl<string>('', [Validators.required]),
        filter: new FormControl<string | null>('', [
            Validators.required,
            Validators.minLength(1),
        ]),
    });

    filters = this.filterService.getFilters().result;
    filterOptions = computed(() =>
        this.filters().data?.map((f) => ({ label: f.name, value: f.id }))
    );

    selectedFile: File | null = null;

    constructor(
        private fb: FormBuilder,
        private filterService: FilterService,
        private jobsService: JobsService
    ) {}

    onFileSelected(event: Event) {
        const file = (event.target as HTMLInputElement).files?.[0];
        if (file) {
            this.selectedFile = file;
            this.checkValidity();
        }
    }

    private checkValidity(): boolean {
        const isValid = this.selectedFile !== null && this.form.valid;
        this.valid.emit(isValid);
        return isValid;
    }

    ngOnInit() {
        this.form.valueChanges.subscribe(() => {
            const valid = this.checkValidity();
            if (valid) {
                this.data.emit({
                    name: this.form.value.name!,
                    filterId: this.form.value.filter!,
                    file: this.selectedFile!,
                });
            }
        });
    }

    touchAll() {
        this.form.markAllAsTouched();
    }

    reset() {
        this.form.reset();
        this.selectedFile = null;
    }
}
