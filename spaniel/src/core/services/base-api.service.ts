import { AxiosInstance } from 'axios';

export interface Headers {
  [key: string]: any;
}

export abstract class BaseApiService<U = any> {
  protected _headers: Headers = {};

  protected _requestInterceptorsIds: number[] = [];
  protected _responseInterceptorsIds: number[] = [];

  constructor (protected axios: AxiosInstance) { }

  public get<T = U>(url: string, headers: Headers = {}) {
    return this.axios.get<T>(`${url}`, {
      headers: this._getHeaders(headers)
    });
  }

  public getAsFile(url: string, headers: Headers = {}) {
    return this.axios.get(`${url}`, {
      headers: this._getHeaders(headers),
      responseType: 'blob'
    });
  }

  public post<T = U>(url: string, data: any, headers: Headers = {}) {
    return this.axios.post<T>(`${url}`, data, {
      headers: this._getHeaders(headers)
    });
  }

  public put<T = U>(url: string, data: any, headers: Headers = {}) {
    return this.axios.put<T>(`${url}`, data, {
      headers: this._getHeaders(headers)
    });
  }

  public delete<T = U>(url: string, headers: Headers = {}) {
    return this.axios.delete<T>(`${url}`, {
      headers: this._getHeaders(headers)
    });
  }

  protected _getHeaders(headers: Headers = {}) {
    return {
      ...this._headers,
      ...headers
    } as Headers;
  }

  protected _setHeader<T = any>(key: string, value: T) {
    this._headers = {
      ...this._headers,
      [key]: value
    };
  }

  protected _ejectInterceptors() {
    this._requestInterceptorsIds.forEach(id => this.axios.interceptors.request.eject(id));
    this._responseInterceptorsIds.forEach(id => this.axios.interceptors.response.eject(id));
  }
}
