import { Component, computed, inject, signal } from '@angular/core';
import { LayoutComponent } from '../../../components/layout/layout.component';
import { ImagesService } from '../../../services/images.service';
import { JsonPipe } from '@angular/common';
import enviroment from '../../../../enviroments/enviroment';
import { ZipService } from '../../../services/zip.service';
import { ActionsComponent } from '../../../components/actions/actions.component';
import { Flowbite } from '../../../util/flowbiteFix';

@Component({
  selector: 'app-found-images-screen',
  standalone: true,
  imports: [LayoutComponent, JsonPipe, ActionsComponent],
  templateUrl: './found-images-screen.component.html',
  styleUrl: './found-images-screen.component.css',
})
export class FoundImagesScreenComponent {
  imagesService = inject(ImagesService);
  zipService = inject(ZipService);
  readonly pageSize = 10;
  currentPageNumber = signal(0);

  actions = [{ id: 'exportAll', label: 'Export all' }];

  onActionSelected(action: { id: string }) {
    switch (action.id) {
      case 'exportAll':
        this.exportAll();
        break;
    }
  }

  imagePages = this.imagesService.getImages(this.pageSize).result;
  currentPage = computed(
    () => this.imagePages().data?.pages[this.currentPageNumber()]
  );

  referencesUrls = computed(() =>
    this.currentPage()?.pictures.map(
      (r) => `${enviroment.api}/files/${r.blob_id}`
    )
  );

  exportAll() {
    this.zipService.downloadZip();
  }

  previous() {
    if (this.imagePages().hasPreviousPage) {
      this.currentPageNumber.update((prev) => prev - 1);
      this.imagePages().fetchPreviousPage();
    }
  }

  next() {
    if (this.imagePages().hasNextPage) {
      this.currentPageNumber.update((prev) => prev + 1);
      this.imagePages().fetchNextPage();
    }
  }
}
