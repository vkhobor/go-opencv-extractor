import {
  Component,
  EventEmitter,
  Output,
  computed,
  effect,
  inject,
  signal,
} from '@angular/core';
import {
  FormBuilder,
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { CreateJob } from '../../../../models/CreateJob';
import { FilterService } from '../../../../services/filter.service';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery } from '@ngneat/query';
import { DefaultHttpProxyService } from '../../../../services/http/default-http-proxy.service';
import { Observable } from 'rxjs';
import { Filter } from '../../../../models/Filter';
import enviroment from '../../../../../enviroments/enviroment';
import { UndefinedInitialDataOptions } from '@ngneat/query/lib/query-options';

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

  selectedFilterIdOrUndefined = toSignal(
    this.form.controls.filter.valueChanges
  );
  selectedFilterId = computed(
    () => this.selectedFilterIdOrUndefined() ?? undefined
  );

  selectedFiltersQuery = this.filterService.getFilter(this.selectedFilterId());

  selectedFilterEffect = effect(() =>
    this.selectedFiltersQuery.updateOptions(
      this.filterService.selectedFiltersQueryOptions(this.selectedFilterId())
    )
  );

  selectedFiltersQueryResult = this.selectedFiltersQuery.result;

  selectedFilterPictures = computed(() => {
    if (!this.selectedFilterId()) {
      return null;
    }

    return this.selectedFiltersQueryResult().data?.filter_images.map((id) => ({
      url: `${enviroment.api}/files/${id.blob_id}`,
    }));
  });

  constructor(private fb: FormBuilder, private filterService: FilterService) {}

  public get searchQuery() {
    return this.form.get('searchQuery');
  }
  public get limit() {
    return this.form.get('limit');
  }

  touchAll() {
    this.form.markAllAsTouched();
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
