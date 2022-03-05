import {
  Enums,
  Statistics,
  StatisticsDto,
} from 'models';

import { ShepherdApiService } from './shepherd-api.service';

export class GeneralService {
  private static URLS = {
    enums: '/getEnums',
    statistics: '/getStatistics'
  };

  public constructor(private apiClient: ShepherdApiService) {}

  public getEnums() {
    const url = GeneralService.URLS.enums;

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
}

