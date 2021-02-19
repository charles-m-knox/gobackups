import { Component, OnInit } from '@angular/core';
import { DataService } from '../data.service';
import { HttpService } from '../http.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  public chartWidth = window.innerWidth * 0.75;
  public chartOptions = {
    responsive: false,
    animation: {
      duration: 0 // prevents the chart from animating every time we update data
    }
  };

  constructor(public http: HttpService, public data: DataService) {
    this.http.getApiData('http://127.0.0.1:12403/api/logs?t=chart', (apiData) => {
      if (apiData) {
        this.data.logData = apiData;
      }
    });
    this.http.getApiData('http://127.0.0.1:12403/api/logs', (apiData) => {
      if (apiData) {
        this.data.logEntries = apiData;
      }
    });
  }

  ngOnInit(): void {
  }

}
