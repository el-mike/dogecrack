import {
  PitbullJobDto,
  mapPitbullJob,
  PitbullJobsFilters,
  ListRequest,
  ListResponse,
  ListResponseResult,
  PitbullJob,
  RunPitbullJobPayload,
} from 'models';

import { ShepherdApiService } from './shepherd-api.service';

export class PitbullJobService {
  private static URLS = {
    jobs: '/getJobs',
    crack: '/crack',
  };

  public constructor(private apiClient: ShepherdApiService) {}

  public getJobs(request: ListRequest<PitbullJobsFilters>) {
    const url = this.apiClient.buildUrl(PitbullJobService.URLS.jobs, request);

    return this.apiClient
      .get<ListResponse<PitbullJobDto>>(url)
      .then(response => ({
        entities: (response.data.data || []).map(job => mapPitbullJob(job)),
        page: response.data.page,
        totalCount: response.data.totalCount,
      } as ListResponseResult<PitbullJob>));
  }

  public runJob(payload: RunPitbullJobPayload) {
    const url = this.apiClient.buildUrl(PitbullJobService.URLS.crack);

    return this.apiClient
      .post<PitbullJobDto>(url, payload)
      .then(response => mapPitbullJob(response.data));
  }
}

