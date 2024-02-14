import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class AddNewService {
  constructor() {}

  modalVisible = false;

  showModal() {
    this.modalVisible = true;
  }
  hideModal() {
    this.modalVisible = false;
  }
}
