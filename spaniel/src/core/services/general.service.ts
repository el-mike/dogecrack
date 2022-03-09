import {
  Enums,
  Statistics,
  StatisticsDto,
  Settings,
  SettingsDto,
} from 'models';

import { ShepherdApiService } from './shepherd-api.service';

export class GeneralService {
  private static URLS = {
    enums: '/getEnums',
    statistics: '/getStatistics',
    settings: '/getSettings',
    updateSettings: '/updateSettings',
  };

  public constructor(private apiClient: ShepherdApiService) {}

  public getEnums() {
    const url = this.apiClient.buildUrl(GeneralService.URLS.enums);

    return this.apiClient
      .get<Enums>(url)
      .then(response => response.data);
  }
  public getStatistics() {
    const url = this.apiClient.buildUrl(GeneralService.URLS.statistics);

    return this.apiClient
      .get<StatisticsDto>(url)
      .then(response => response.data as Statistics);
  }

  public getSettings() {
    const url = this.apiClient.buildUrl(GeneralService.URLS.settings);

    return this.apiClient
      .get<SettingsDto>(url)
      .then(response => response.data as Settings);
  }

  public updateSettings(payload: Settings) {
    const url = this.apiClient.buildUrl(GeneralService.URLS.updateSettings);

    return this.apiClient
      .patch<SettingsDto>(url, payload)
      .then(response => response.data as Settings)
  }
}

