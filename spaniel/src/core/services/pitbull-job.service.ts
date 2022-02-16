import {
  PitbullJobDto,
  mapPitbullJob,
} from 'models';

import { ShepherdApiService } from './shepherd-api.service';

export class PitbullJobService {
  private static URLS = {
    jobs: '/getJobs',
  };

  public constructor(private apiClient: ShepherdApiService) {}

  public getJobs() {
    const url = PitbullJobService.URLS.jobs;

    return this.apiClient
      .get<PitbullJobDto[]>(url)
      .then(response => (response.data || []).map(job => mapPitbullJob(job)));
  }
}

