import {
  User,
  UserCredentials,
} from 'models';

import { ShepherdApiService } from './shepherd-api.service';

export class AuthService {
  private static URLS = {
    me: '/getMe',
    login: '/login',
    logout: '/logout'
  };

  public constructor(private apiClient: ShepherdApiService) {}

  public getMe() {
    const url = AuthService.URLS.me;

    return this.apiClient
      .get<User>(url)
      .then(response => response.data);
  }

  public login(creds: UserCredentials) {
    const url = AuthService.URLS.login;
  
    return this.apiClient.post(url, creds);
  }

  public logout() {
    const url = AuthService.URLS.logout;
  
    return this.apiClient.get(url);
  }
}

