import { Component } from '@angular/core';
import { DataService } from './data.service';
import { HelpersService } from './helpers.service';
import { HttpService } from './http.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  public navbarCollapseShow = false;

  constructor(public http: HttpService, public data: DataService, public helpers: HelpersService) { }

  /**
   * For certain mobile/desktop breakpoints, this function gets called. It
   * simply flips the boolean value that is bound to showing/hiding the
   * navbar collapse button in the top right of the header.
   * @returns void
   */
  public toggleNavCollapse(): void {
    this.navbarCollapseShow = !this.navbarCollapseShow;
  }
}
