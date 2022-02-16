import {
  Enums,
} from 'models';

import { ShepherdApiService } from './shepherd-api.service';

export class GeneralService {
  private static URLS = {
    jobs: '/getEnums',
  };

  public constructor(private apiClient: ShepherdApiService) {}

  public getEnums() {
    const url = GeneralService.URLS.jobs;

    return this.apiClient
      .get<Enums>(url)
      .then(response => response.data);
  }
}

