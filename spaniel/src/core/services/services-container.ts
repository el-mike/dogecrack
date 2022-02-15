import axios from 'axios';

import { ShepherdApiService } from './shepherd-api.service';
import { ConfigService } from './config.service';
import { AuthService } from './auth.service';
import { PitbullJobService } from './pitbull-job.service';

export type UnauthorizedResponseCallback = () => void;

/**
 * This file acts as a simple dependencies container, that allows us to facilitate easy service creation
 * and distribution. Services are stateless (except for BaseApiService, which has globally-applied config),
 * therefore we can safely reuse their instances.
 */

export const configService = ConfigService.getInstance();

export const shepherdApiService = new ShepherdApiService(axios, configService);

export const authService = new AuthService(shepherdApiService);
export const pitbullJobService = new PitbullJobService(shepherdApiService);
