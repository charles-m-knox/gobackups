import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DataService } from './data.service';

@Injectable({
  providedIn: 'root'
})
export class HttpService {

  constructor(public http: HttpClient, public data: DataService) { }

  /**
   * Logs an error message to console.error, and also returns the message.
   * @param {string} source Source of what called the logging function, could be anything that helps you trace
   * @param {string} message The message to log
   * @returns string
   */
  public logError(source: string, message: string): string {
    const output = `${source}: ${message}`;
    console.error(output);
    return output;
  }

  /**
   * getApiData is a wrapper around the observable "get" function. Pass in a
   * function that will get executed after the get request.
   * @param {string} url URL to request
   * @param {any} func Arrow function to be executed after the request
   *
   * **Example:**
   * ```
   * func: any = (data) => { this.data.foo = data.bar; }
   * ```
   * @param {any} options? HTTP options to include when making the request
   * @returns void
   */
  public getApiData(url: string, func: any, options?: any): void {
    this.http.get(url, options)
      .subscribe((data: any) => {
        func(data);
      }, error => {
        this.logError('getApiData', `API get error ${error.status}: ${error.message}`);
      });
  }
  /**
   * postApiData is a wrapper around the observable "post" function. Pass in a
   * function that will get executed after the post request.
   * @param {string} url URL to request
   * @param {any} content Content of the POST body
   * @param {any} func Arrow function to be executed after the request
   *
   * **Example:**
   * ```
   * func: any = (data) => { this.data.foo = data.bar; }
   * ```
   * @param {any} options? HTTP options to include when making the request
   * @returns void
   */
  public postApiData(url: string, content: any, func: any, options?: any): void {
    this.http.post(url, content, options)
      .subscribe((data: any) => {
        func(data);
      }, error => {
        this.logError('postApiData', `API post error ${error.status}: ${error.message}`);
      });
  }

  /**
   * putApiData is a wrapper around the observable "put" function. Pass in a
   * function that will get executed after the put request.
   * @param {string} url URL to request
   * @param {any} content Content of the PUT body
   * @param {any} func Arrow function to be executed after the request
   *
   * **Example:**
   * ```
   * func: any = (data) => { this.data.foo = data.bar; }
   * ```
   * @param {any} options? HTTP options to include when making the request
   * @returns void
   */
  public putApiData(url: string, content: any, func: any, options?: any): void {
    this.http.put(url, content, options)
      .subscribe((data: any) => {
        func(data);
      }, error => {
        this.logError('putApiData', `API put error ${error.status}: ${error.message}`);
      });
  }

  /**
   * deleteApiData is a wrapper around the observable "delete" function. Pass in a
   * function that will get executed after the delete request.
   * @param {string} url URL to request
   * @param {any} func Arrow function to be executed after the request
   *
   * **Example:**
   * ```
   * func: any = (data) => { this.data.foo = data.bar; }
   * ```
   * @param {any} options? HTTP options to include when making the request
   * @returns void
   */
  public deleteApiData(url: string, func: any, options?: any): void {
    this.http.delete(url, options)
      .subscribe((data: any) => {
        func(data);
      }, error => {
        this.logError('deleteApiData', `API delete error ${error.status}: ${error.message}`);
      });
  }
}
