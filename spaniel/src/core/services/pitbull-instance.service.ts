import {
  PitbullInstancesStatisticsDto,
  PitbullInstancesStatistics,
} from 'models';

import { ShepherdApiService } from './shepherd-api.service';

export class PitbullInstanceService {
  private static URLS = {
    statistics: '/getInstancesStatistics',
  };

  public constructor(private apiClient: ShepherdApiService) {}

  public getStatistics() {
    const url = this.apiClient.buildUrl(PitbullInstanceService.URLS.statistics);

    return this.apiClient
      .get<PitbullInstancesStatisticsDto>(url)
      .then(response => response.data as PitbullInstancesStatistics);
  }
}

