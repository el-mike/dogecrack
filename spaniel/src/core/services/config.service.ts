import { Environment } from 'models';

import {
  localConfig,
  prodConfig,
} from 'config/environments';

export class ConfigService {
  private static _instance: ConfigService;
  private _currentEnv: Environment;

  private constructor () {
    /**
     * REACT_APP_ENV is a custom env variable, set directly in
     * app scripts in package.json.
     */
    const nodeEnv = process.env.REACT_APP_ENV;

    this._currentEnv = this._isEnvCorrect(nodeEnv)
      ? nodeEnv
      : Environment.LOCAL;
  }

  public static getInstance() {
    if (!ConfigService._instance) {
      ConfigService._instance = new ConfigService();
    }

    return ConfigService._instance;
  }

  public getAppConfig() {
    if (this._currentEnv === Environment.LOCAL) {
      return localConfig;
    }


    if (this._currentEnv === Environment.PROD) {
      return prodConfig;
    }

    return localConfig;
  }

  public isProd() {
    return this._currentEnv === Environment.PROD;
  }

  public getCurrentEnv() {
    return this._currentEnv;
  }

  private _isEnvCorrect(env: string | undefined): env is Environment {
    return (
      !!env && Object.values<string>(Environment).includes(env.toLowerCase())
    );
  }
}
