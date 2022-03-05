import { ShepherdApiService } from './shepherd-api.service';

export class PitbullInstanceService {
  private static URLS = {
    statistics: '/getInstancesStatistics',
  };

  public constructor(private apiClient: ShepherdApiService) {}
}

