import { Injectable } from '@angular/core';
import { DataService } from './data.service';

@Injectable({
  providedIn: 'root'
})
export class HelpersService {

  constructor(public data: DataService) { }

  /**
   * `setPage` is the function that will allow any component to change the current tab/page/view.
   * @param {string} page The name of the page to change to.
   * @returns void
   */
  public setPage(page: string): void {
    this.data.currentPage = page;
  }
}
