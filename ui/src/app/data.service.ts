import { Injectable } from '@angular/core';

/**
 * Represents a page that users can navigate to, either via
 * a URL route, the navbar, or from other component interactions.
 *
 * @interface
 */
export interface Page {
  page: string;
  caption: string;
}

@Injectable({
  providedIn: 'root'
})
export class DataService {

  constructor() { }

  // currentPage is set whenever the user taps/clicks on a tab or navigates
  // to a page via the routing module. Its value is important because it is
  // visually bound to highlighting the currently viewed tab
  public currentPage = 'home';

  // pages is a list of all tabs that are navigable by the user, accessible
  // via the routing module and also by clicking tabs
  public pages: Page[] = [
    { page: 'home', caption: 'Home' },
  ];

  public logData: any = {
    series: [],
    labels: [],
  };

  public logEntries: any = [];
}
