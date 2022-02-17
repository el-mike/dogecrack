import axios, {
  AxiosRequestConfig,
} from 'axios';

import { Dictionary } from 'models';

import { toQueryParams } from 'core/utils';

import { BaseApiService } from './base-api.service';

import { ConfigService } from './config.service';

interface CustomRequestConfig extends AxiosRequestConfig {
  retry?: boolean;
}

const TOKEN_HEADER = 'Authorization';
const CONTENT_TYPE_HEADER = 'Content-Type';

const TOKEN_PREFIX = 'Bearer';

const UNAUTHORIZED_STATUS = 401;

/**
 * ShepherdApiService - service that provides logic for communicating with ShepherdAPI.
 *
 * ShepherdApiService uses Axios under the hood, and provides a way to set up custom interceptors.
 */
export class ShepherdApiService extends BaseApiService {
  public constructor (
    axiosProvider: typeof axios,
    configService: ConfigService,
  ) {
    super(axiosProvider.create({ baseURL: `${configService.getAppConfig().apiUrl}` }));
  }

  public toQueryParams(request: Dictionary) {
    return toQueryParams(request);
  }

  public buildUrl(baseUrl: string, request: Dictionary = {}) {
    return `${baseUrl}${this.toQueryParams(request)}`;
  }

  /**
   * setInterceptors - sets Axios interceptors for handling authorization header and token refreshing.
   * It allows to specify callback function for refresh token, which will be fired upon successful call to
   * AuthService.
   */
  public setInterceptors(refreshCallback?: () => void) {
    this._ejectInterceptors();

    this._requestInterceptorsIds.push(
      this.axios.interceptors.request.use(
        async config => {
          return this._attachToken(config, '');
        },
        error => Promise.reject(error)
      )
    );

    this._responseInterceptorsIds.push(
      this.axios.interceptors.response.use(
        response => response,
        async error => {
          const request = error.config as CustomRequestConfig;

          if (error.response?.status === UNAUTHORIZED_STATUS && !request.retry) {
            request.retry = true;

            try {
              this.setInterceptors(refreshCallback);
              refreshCallback?.();

              return this.axios(request);
            } catch (error: unknown) {
              throw error;
            }
          } else {
            throw error;
          }
        }
      )
    );
  }

  private _attachToken(config: AxiosRequestConfig, token: string) {
    config.headers = config.headers || {};
  
    config.headers[TOKEN_HEADER] = `${TOKEN_PREFIX} ${token}`;
    config.headers[CONTENT_TYPE_HEADER] = 'application/json';

    return config;
  }
}
