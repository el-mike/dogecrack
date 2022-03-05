import {
  CrackJobDto,
  mapCrackJob,
  CrackJobsFilters,
  ListRequest,
  ListResponse,
  ListResponseResult,
  CrackJob,
  RunCrackJobPayload,
  CrackJobsStatisticsDto,
  CrackJobsStatistics,
} from 'models';

import { ShepherdApiService } from './shepherd-api.service';

export class CrackJobService {
  private static URLS = {
    jobs: '/getJobs',
    statistics: '/getJobsStatistics',
    crack: '/crack',
  };

  public constructor(private apiClient: ShepherdApiService) {}

  public getJobs(request: ListRequest<CrackJobsFilters>) {
    const url = this.apiClient.buildUrl(CrackJobService.URLS.jobs, request);

    return this.apiClient
      .get<ListResponse<CrackJobDto>>(url)
      .then(response => ({
        entities: (response.data.data || []).map(job => mapCrackJob(job)),
        page: response.data.page,
        totalCount: response.data.totalCount,
      } as ListResponseResult<CrackJob>));
  }

  public getStatistics() {
    const url = this.apiClient.buildUrl(CrackJobService.URLS.statistics);

    return this.apiClient
      .get<CrackJobsStatisticsDto>(url)
      .then(response => response.data as CrackJobsStatistics);
  }

  public runJob(payload: RunCrackJobPayload) {
    const url = this.apiClient.buildUrl(CrackJobService.URLS.crack);

    return this.apiClient
      .post<CrackJobDto>(url, payload)
      .then(response => mapCrackJob(response.data));
  }
}

