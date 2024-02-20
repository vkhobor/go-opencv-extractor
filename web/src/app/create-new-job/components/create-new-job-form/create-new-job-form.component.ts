import { Component, EventEmitter, Output } from '@angular/core';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { CreateJob } from '../../../models/CreateJob';

@Component({
  selector: 'app-create-new-job-form',
  standalone: true,
  imports: [FormsModule, ReactiveFormsModule],
  templateUrl: './create-new-job-form.component.html',
  styleUrl: './create-new-job-form.component.css',
})
export class CreateNewJobFormComponent {
  @Output() valid = new EventEmitter<boolean>();
  @Output() data = new EventEmitter<CreateJob | undefined>();

  form = this.fb.group(
    {
      searchQuery: new FormControl<string | null>(null, [
        Validators.required,
        Validators.minLength(5),
      ]),
      limit: new FormControl<number | null>(null, [
        Validators.required,
        Validators.min(1),
        Validators.max(1000),
      ]),
    },
    { updateOn: 'change' }
  );

  constructor(private fb: FormBuilder) {}

  reset() {
    this.form.reset();
  }

  public get searchQuery() {
    return this.form.get('searchQuery');
  }
  public get limit() {
    return this.form.get('limit');
  }

  ngOnInit() {
    this.form.valueChanges.subscribe((data) => {
      this.valid.emit(this.form.valid);

      if (this.form.valid) {
        this.data.emit({
          search_query: data.searchQuery!,
          limit: data.limit!,
        });
      } else {
        this.data.emit(undefined);
      }
    });
  }
}
